terraform {
  required_providers {
    instatus = {
      # Using a local dev override for this
      source = "ashleyjackson/instatus"
      # version = "~> 0.1"
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

  email          = "ashley@myaffiliates.com"
  workspace_slug = each.value
  name           = "Instatus Page for ${each.value}"
}

