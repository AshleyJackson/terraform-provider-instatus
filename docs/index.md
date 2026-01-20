---
page_title: "Instatus Provider"
subcategory: ""
description: |-
  Terraform provider for managing Instatus status pages and components.
---

# Instatus Provider

The Instatus provider allows you to manage [Instatus](https://instatus.com) status pages and components using Terraform.

## Example Usage

```terraform
terraform {
  required_providers {
    instatus = {
      source  = "ashleyjackson/instatus"
      version = "~> 0.1"
    }
  }
}

provider "instatus" {
  api_key = var.instatus_api_key
}
```

## Authentication

The provider requires an API key which can be obtained from your Instatus dashboard.

## Schema

### Required

* `api_key` - (String, Sensitive) Instatus API key for authentication.

## Resources

- [instatus_component](resources/component) - Manage status page components
