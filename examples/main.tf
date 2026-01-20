terraform {
  required_providers {
    instatus = {
      # Using a local dev override for this
      source = "ashleyjackson/instatus"
      version = "0.1.4"
    }
  }
}

provider "instatus" {
  api_key = var.instatus_api_key
}

variable "programs" {
  type    = list(string)
  default = ["program1", "program2"]
}

resource "instatus_page" "pages" {
  for_each = toset(var.programs)

  email          = var.email
  workspace_slug = each.value
  name           = "Instatus Page for ${each.value}"
}

