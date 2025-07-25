
locals {
  current_datetime = timestamp()
  azure_region     = "northcentralus"
  tags = {
    Environment  = "Live",
    Application  = "Azure Tagger",
    CreatedBy    = var.tag_creator_name,
    ManagedBy    = "Azure Team",
    CreationDate = local.current_datetime
    deployment   = "terraform",
  }

  owners = [
    data.azurerm_client_config.current.object_id,
    # data.azuread_group.admins.object_id,
    # [for role in data.azuread_directory_roles.global_admin.roles : role.object_id if role.display_name == "Global Administrator"]
  ]
}

# Define the naming convention for resources
module "naming" {
  source = "Azure/naming/azurerm"
  prefix = ["acme"]
  suffix = [
    "azure",
    "tagger",
  ]
}

# Resource Group
resource "azurerm_resource_group" "azure_tagger" {
  name     = module.naming.resource_group.name_unique
  location = var.location
  tags     = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

resource "azurerm_log_analytics_workspace" "azure_tagger_law" {
  name                = module.naming.log_analytics_workspace.name_unique
  resource_group_name = azurerm_resource_group.azure_tagger.name
  location            = azurerm_resource_group.azure_tagger.location
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags                = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

resource "azurerm_application_insights" "azure_tagger_ai" {
  name                = module.naming.application_insights.name_unique
  resource_group_name = azurerm_resource_group.azure_tagger.name
  location            = azurerm_resource_group.azure_tagger.location
  workspace_id        = azurerm_log_analytics_workspace.azure_tagger_law.id
  application_type    = "web"
  tags                = local.tags
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

# User Assigned Managed Identity
resource "azurerm_user_assigned_identity" "azure_tagger" {
  name                = module.naming.user_assigned_identity.name_unique
  location            = azurerm_resource_group.azure_tagger.location
  resource_group_name = azurerm_resource_group.azure_tagger.name
  tags                = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

# Granting the UAMI permissions to update tags on resources in subscription
resource "azurerm_role_assignment" "uami_resource_tagging" {
  principal_id         = azurerm_user_assigned_identity.azure_tagger.principal_id
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Tag Contributor"
}

# Granting the UAMI permissions to read resources in subscription (to read tags)
resource "azurerm_role_assignment" "uami_resource_reading" {
  principal_id         = azurerm_user_assigned_identity.azure_tagger.principal_id
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Reader"
}

# Storage Account for Function App
resource "azurerm_storage_account" "azure_tagger" {
  name                     = module.naming.storage_account.name_unique
  resource_group_name      = azurerm_resource_group.azure_tagger.name
  location                 = azurerm_resource_group.azure_tagger.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  tags                     = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

# Storage Container for the Function App
resource "azurerm_storage_container" "azure_tagger" {
  name                  = "acme-function-releases"
  storage_account_name  = azurerm_storage_account.azure_tagger.name
  container_access_type = "private"
}

# Assign Contributor role to Function Apps' assigned managed identity
resource "azurerm_role_assignment" "rbac_storage_blob_contributor_azure_tagger" {
  scope                = azurerm_storage_account.azure_tagger.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.azure_tagger.principal_id
}

# Application Service Plan
resource "azurerm_service_plan" "azure_tagger" {
  name                = module.naming.app_service_plan.name_unique
  location            = azurerm_resource_group.azure_tagger.location
  resource_group_name = azurerm_resource_group.azure_tagger.name
  os_type             = "Linux"
  sku_name            = "Y1"
  tags                = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

# Function App
resource "azurerm_linux_function_app" "azure_tagger" {
  name                 = module.naming.function_app.name_unique
  location             = azurerm_resource_group.azure_tagger.location
  resource_group_name  = azurerm_resource_group.azure_tagger.name
  service_plan_id      = azurerm_service_plan.azure_tagger.id
  storage_account_name = azurerm_storage_account.azure_tagger.name
  # storage_account_access_key = azurerm_storage_account.azure_tagger.primary_access_key
  storage_uses_managed_identity = true
  # content_share_force_disabled = true
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.azure_tagger.id]
  }
  https_only = true

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "custom"
    // AZURE_CLIENT_ID is picked up by azidentity.NewDefaultAzureCredential(nil)
    AZURE_CLIENT_ID       = azurerm_user_assigned_identity.azure_tagger.client_id
    AZURE_SUBSCRIPTION_ID = var.az_subscription_id
    AZURE_TAGGER_PREFIX   = var.tag_prefix
    # WEBSITE_CONTENTSHARE                     = azurerm_storage_share.azure_tagger.name
    WEBSITE_RUN_FROM_PACKAGE_BLOB_MI_RESOURCE_ID = azurerm_user_assigned_identity.azure_tagger.id
    WEBSITE_AUTH_AAD_ALLOWED_TENANTS             = data.azurerm_subscription.primary.tenant_id
  }

  auth_settings_v2 {
    auth_enabled             = true
    default_provider         = "aad"
    excluded_paths           = []
    forward_proxy_convention = "NoProxy"
    http_route_api_prefix    = "/.auth"
    require_authentication   = true
    require_https            = true
    runtime_version          = "~1"
    unauthenticated_action   = "Return401"
    active_directory_v2 {
      allowed_applications            = [data.azuread_service_principal.eventgrid.client_id]
      allowed_audiences               = [format("api://%s", azuread_application.appreg_azure_tagger.client_id)]
      allowed_groups                  = []
      allowed_identities              = []
      client_id                       = azuread_application.appreg_azure_tagger.client_id
      jwt_allowed_client_applications = []
      jwt_allowed_groups              = []
      login_parameters                = {}
      tenant_auth_endpoint            = format("https://sts.windows.net/%s/v2.0", data.azurerm_subscription.primary.tenant_id)
      www_authentication_disabled     = false
    }
    login {
      allowed_external_redirect_urls    = []
      cookie_expiration_convention      = "FixedTime"
      cookie_expiration_time            = "08:00:00"
      logout_endpoint                   = "/.auth/logout"
      nonce_expiration_time             = "00:05:00"
      preserve_url_fragments_for_logins = false
      token_refresh_extension_time      = 72
      token_store_enabled               = false
      validate_nonce                    = true
    }
  }
  site_config {
    minimum_tls_version                    = "1.2"
    application_insights_connection_string = azurerm_application_insights.azure_tagger_ai.connection_string
    application_stack {
      use_custom_runtime = true
    }
    scm_minimum_tls_version = "1.2"
  }
  tags = local.tags

  lifecycle {
    ignore_changes = [
      app_settings["WEBSITE_RUN_FROM_PACKAGE"],
      tags
    ]
  }
}

# App Registration
resource "azuread_application" "appreg_azure_tagger" {
  display_name     = "acme_appreg_azure_tagger"
  sign_in_audience = "AzureADMyOrg"
  owners           = local.owners

  api {
    mapped_claims_enabled          = false
    requested_access_token_version = 2
    oauth2_permission_scope {
      admin_consent_description  = "Allow users to call Azure Tagger"
      admin_consent_display_name = "Allow users to call Azure Tagger"
      enabled                    = true
      id                         = var.appreg_azure_tagger_uuid
      type                       = "User"
      value                      = "user_impersonation"
    }

  }

  timeouts {}

  web {

    implicit_grant {
      access_token_issuance_enabled = false
      id_token_issuance_enabled     = false
    }
  }
  lifecycle {
    ignore_changes = [
      identifier_uris,
    ]
  }
}

resource "azuread_application_identifier_uri" "appreg_azure_tagger" {
  application_id = azuread_application.appreg_azure_tagger.id
  identifier_uri = format("api://%s", azuread_application.appreg_azure_tagger.client_id)
}

resource "azuread_service_principal" "azure_tagger" {
  client_id = azuread_application.appreg_azure_tagger.client_id
  owners    = local.owners
  feature_tags {
    custom_single_sign_on = false
    enterprise            = true
    hide                  = true
  }
  timeouts {}
}

# Assign Microsoft Graph API permissions (Directory.Read.All) to azure_tagger UAMI
resource "azuread_app_role_assignment" "graph_directory_read_all" {
  principal_object_id = azurerm_user_assigned_identity.azure_tagger.principal_id
  app_role_id         = data.azuread_service_principal.microsoft_graph.app_role_ids["Directory.Read.All"]
  resource_object_id  = data.azuread_service_principal.microsoft_graph.object_id
}

# Event Grid System Topic where events from Azure Subscription are published
resource "azurerm_eventgrid_system_topic" "azure_tagger" {
  # name                   = "azure-tagger-topic"
  name                   = module.naming.eventgrid_topic.name_unique
  resource_group_name    = azurerm_resource_group.azure_tagger.name
  location               = "Global"
  source_arm_resource_id = data.azurerm_subscription.primary.id
  topic_type             = "Microsoft.Resources.Subscriptions"
  tags                   = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

# Event Grid Subscription to trigger the Function App
resource "azurerm_eventgrid_event_subscription" "azure_tagger" {
  # name  = "azure-tagger-event-subscription"
  name  = module.naming.eventgrid_event_subscription.name_unique
  scope = data.azurerm_subscription.primary.id
  webhook_endpoint {
    # url                               = format("https://%s/api/SendGridEvents", azurerm_linux_function_app.azure_tagger.default_hostname)
    url                               = format("https://%s.azurewebsites.net/api/SendGridEvents", azurerm_linux_function_app.azure_tagger.name)
    max_events_per_batch              = 1
    preferred_batch_size_in_kilobytes = 64
    active_directory_tenant_id        = data.azurerm_subscription.primary.tenant_id
    active_directory_app_id_or_uri    = azuread_application.appreg_azure_tagger.client_id
  }
  included_event_types = ["Microsoft.Resources.ResourceWriteSuccess"]

  retry_policy {
    event_time_to_live    = 1440
    max_delivery_attempts = 30
  }
}

resource "azurerm_monitor_action_group" "action_group" {
  name                = "Application Insights Smart Detection"
  resource_group_name = azurerm_resource_group.azure_tagger.name
  short_name          = "SmartDetect"
  tags                = local.tags

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
