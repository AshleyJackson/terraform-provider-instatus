# Terraform Provider for Instatus

Manage Instatus status pages and components with Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

```shell
git clone https://github.com/ashleyjackson/terraform-provider-instatus
cd terraform-provider-instatus
make install
```

## Using the Provider

```hcl
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

resource "instatus_page" "example" {
  email          = "admin@example.com"
  name           = "Example Status"
  workspace_slug = "example"
}
```

## DO NOT USE IN YOUR PRODUCTION ENV
Please note that due to a lack of singular Get Status Pages and Workspaces API, this uses a custom Internal integration that performs this lookup and provides a quick result.

## Documentation

Full documentation is available in the [docs](./docs) directory or on the [Terraform Registry](https://registry.terraform.io/providers/ashleyjackson/instatus/latest/docs).

## Development

### Local Testing

1. Build and install locally:
```shell
make install
```

2. Create `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "ashleyjackson/instatus" = "/home/ashley/go/bin"
  }
  direct {}
}
```

3. Run terraform in examples directory:
```shell
cd examples
terraform init
terraform plan
```

### Running Tests

```shell
make test
make testacc
```

## License

Mozilla Public License v2.0
