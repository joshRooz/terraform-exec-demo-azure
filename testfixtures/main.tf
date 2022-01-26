provider "azurerm" {
  features {}
}

variable "resource_group_name" { type = string }
variable "location" { type = string }

resource "azurerm_resource_group" "test" {
  name     = var.resource_group_name
  location = var.location
}