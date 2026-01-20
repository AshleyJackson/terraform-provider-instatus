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

resource "instatus_page" "pages" {
  name           = "MyAffiliates Status Page2"
  email          = "ashley@myaffiliates.com"
  workspace_slug = "myaffiliates-status2"
}

