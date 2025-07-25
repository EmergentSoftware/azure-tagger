
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
  default = "North Central US"
}

variable "client_name_prefix" {
  default = "acme"
}

# unique id
variable "appreg_azure_tagger_uuid" {
  default = "02c0cc5e-45e8-41d3-b71a-3530923b1fae"
}

variable "tag_prefix" {
  default = ""
}

variable "tag_creator_name" {
  default = "Josh Dinndorf"
}


variable "az_tenant_id" {
  type = string
}

variable "az_subscription_id" {
  type = string
}