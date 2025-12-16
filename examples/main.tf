terraform {
  required_providers {
    instatus = {
      # Using a local dev override for this
      source = "ashleyjackson/instatus"
      # version = "~> 0.1"
    }
  }
}

variable "instatus_api_key" {
  description = "Instatus API key"
  type        = string
  sensitive   = true
}


provider "instatus" {
  api_key = var.instatus_api_key
}

# resource "instatus_page" "pages" {
#   name      = ""
#   email     = ""
#   subdomain = ""
# }

