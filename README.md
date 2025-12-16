# Terraform Provider for Instatus

A Terraform provider for managing [Instatus](https://instatus.com) status page components.

## Features

- ✅ Create components
- ✅ Read component details
- ✅ Update components
- ✅ Delete components
- ✅ Support for nested components (groups)
- ✅ Support for component status, ordering, and archiving

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- Instatus API Key

## Installation

### For Development

1. Clone the repository:
```bash
git clone https://github.com/ashleyjackson/terraform-provider-instatus.git
cd terraform-provider-instatus
```

2. Build the provider:
```bash
go build -o terraform-provider-instatus
```

3. Install the provider locally:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ashleyjackson/instatus/0.1.0/darwin_arm64
cp terraform-provider-instatus ~/.terraform.d/plugins/registry.terraform.io/ashleyjackson/instatus/0.1.0/darwin_arm64/
```

Note: Replace `darwin_arm64` with your OS/architecture:
- macOS (Intel): `darwin_amd64`
- macOS (Apple Silicon): `darwin_arm64`
- Linux: `linux_amd64`
- Windows: `windows_amd64`

## Usage

### Provider Configuration

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
  api_key = var.instatus_api_key  # or set INSTATUS_API_KEY env var
}
```

### Creating a Simple Component

```hcl
resource "instatus_component" "api" {
  name        = "API Service"
  description = "Our main API service"
  status      = "OPERATIONAL"
  show_uptime = true
}
```

### Creating a Nested Component (Group with Children)

```hcl
# Parent group
resource "instatus_component" "web_services" {
  name        = "Web Services"
  description = "All web-related services"
  status      = "OPERATIONAL"
  show_uptime = true
}

# Child component
resource "instatus_component" "website" {
  name        = "Website"
  description = "Main website"
  status      = "OPERATIONAL"
  show_uptime = true
  grouped     = true
  group_id    = instatus_component.web_services.id
}
```

### Component Status Values

- `OPERATIONAL` - Component is working normally
- `UNDERMAINTENANCE` - Component is under maintenance
- `DEGRADEDPERFORMANCE` - Component has degraded performance
- `PARTIALOUTAGE` - Component has partial outage
- `MAJOROUTAGE` - Component has major outage

## Resource: instatus_component

### Arguments

- `name` (Required) - The name of the component
- `description` (Optional) - The description of the component. Default: ""
- `status` (Optional) - The status of the component. Default: "OPERATIONAL"
- `show_uptime` (Optional) - Whether to show uptime for this component. Default: true
- `order` (Optional, Computed) - The order of the component. If not set, will be managed in the Instatus UI. Terraform will read the value from API but won't change it unless explicitly set.
- `grouped` (Optional) - Whether this component belongs to a group. Default: false
- `group_id` (Optional) - The ID of the parent group (required if grouped is true)
- `archived` (Optional) - Whether the component is archived. Default: false

### Attributes

- `id` - The ID of the component
- `unique_email` - The unique email address for this component (for automation)

## Getting API Credentials

1. Log in to your [Instatus dashboard](https://instatus.com)
2. Navigate to **Settings** → **API**
3. Generate an API key

Set as environment variables:

```bash
export INSTATUS_API_KEY="your-api-key"
```

## Development

### Building the Provider

```bash
go build -o terraform-provider-instatus
```

### Running Tests

```bash
go test ./...
```

### Initialize Go Modules

```bash
go mod tidy
```

## API Documentation

This provider uses the Instatus API:
- [API Documentation](https://instatus.com/help/api/components)
- API v1 for CREATE and DELETE operations
- API v2 for READ and UPDATE operations

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
