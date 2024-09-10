
variable "azure_subscription_id" {
  type    = string
}

# Get tenant and subscription(Azure formatted) data
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "current" {}

# Retrieve Global Administrator role ID
# data "azurerm_role_definition" "global_admin" {
#   name = "Global Administrator"
# }

# Get the Microsoft Event Grid enterprise application
data "azuread_service_principal" "eventgrid" {
  display_name = "Microsoft.EventGrid"
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
    # data.azurerm_role_definition.global_admin.object_id
  ]
}