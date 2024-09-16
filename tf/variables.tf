
variable "azure_subscription_id" {
  type = string
}

# Get tenant and subscription(Azure formatted) data
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "current" {}

# # Retrieve Global Administrator role ID
# data "azuread_directory_roles" "global_admin" {
#   # display_name = "Global Administrator"
# }
# data "azuread_group" "admins" {
#   display_name = "admins"
# }

# Get the Microsoft Event Grid enterprise application
data "azuread_service_principal" "eventgrid" {
  display_name = "Microsoft.EventGrid"
}

# Get Microsoft Graph Service Principal to pick Directory.Read.All permission from it
data "azuread_service_principal" "microsoft_graph" {
  display_name = "Microsoft Graph"
}

variable "location" {
  default = "Central US"
}

variable "resource_group" {
  default = "rg-azure-tagger"
}

variable "log_analytics_workspace_name" {
  default = "azure-tagger-law"
}

variable "application_insights_name" {
  default = "azure-tagger-appinsights"
}

variable "tagger_storage_account" {
  default = "azuretaggerstorage"
}

variable "azure_tagger_function_name" {
  default = "azure-tagger-function"
}

# unique id
variable "appreg_azure_tagger_uuid" {
  default = "02c0cc5e-45e8-41d3-b71a-3530923b1fae"
}

variable "tag_prefix" {
  default = "AzTagger"
}

variable "tag_creator_name" {
  default = "AzureTaggerTerraform"
}

locals {
  current_datetime = timestamp()
  # azuread_application must have current user as an owners to authenticate
  owners = [
    data.azurerm_client_config.current.object_id,
    # data.azuread_group.admins.object_id,
    # [for role in data.azuread_directory_roles.global_admin.roles : role.object_id if role.display_name == "Global Administrator"]
  ]
}
