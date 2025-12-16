---
page_title: "Instatus Provider"
subcategory: ""
description: |-
  Terraform provider for managing Instatus status page components.
---

# Instatus Provider

The Instatus provider allows you to manage [Instatus](https://instatus.com) status page components using Terraform.

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

The provider requires an Instatus API key. These can be provided in two ways:

### Environment Variables

```bash
export INSTATUS_API_KEY="your-api-key"
```

### Provider Configuration

```terraform
provider "instatus" {
  api_key = "your-api-key"
}
```

## Getting API Credentials

1. Log in to your [Instatus dashboard](https://instatus.com)
2. Navigate to **Settings** → **API**
3. Generate an API key

## Schema

### Required

- `api_key` (String, Sensitive) - Instatus API key for authentication

## Resources

- [instatus_component](resources/component) - Manage status page components
