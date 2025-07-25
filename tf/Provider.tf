terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.31.0"
    }
  }

}

provider "azurerm" {
  subscription_id                 = var.az_subscription_id
  tenant_id                       = var.az_tenant_id

  features {}
}
