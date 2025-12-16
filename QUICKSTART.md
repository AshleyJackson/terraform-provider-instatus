# Quick Start Guide - Instatus Terraform Provider

## What We Built

A fully functional Terraform provider for managing Instatus status page components with CRUD operations and support for nested component hierarchies (up to 3 levels).

## Project Structure

```
terraform-provider-instatus/
├── main.go                  # Provider entry point
├── go.mod                   # Go module dependencies
├── instatus/
│   ├── provider.go          # Provider configuration
│   ├── client.go            # Instatus API client
│   └── resource_component.go # Component resource implementation
├── examples/
│   ├── main.tf              # Example Terraform configuration
│   └── variables.tf         # Variable definitions
├── test/
│   ├── main.tf              # Simple test configuration
│   └── nested_test.tf       # Nested hierarchy test
└── README.md                # Comprehensive documentation
```

## Key Features Implemented

✅ **CREATE** - Create new components with all attributes
✅ **READ** - Read component details and refresh state
✅ **UPDATE** - Update component properties
✅ **DELETE** - Delete components
✅ **Nested Components** - Support for 3-level hierarchies (parent → child → grandchild)
✅ **Environment Variables** - Configure via `INSTATUS_API_KEY`

## API Compatibility

- **v1 API**: Used for CREATE and DELETE operations
- **v2 API**: Used for READ and UPDATE operations (supports hierarchical responses)

## Building the Provider

```bash
cd /terraform-provider-instatus
go build -o terraform-provider-instatus
```

## Installing Locally

```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ashleyjackson/instatus/0.1.0/darwin_arm64
cp terraform-provider-instatus ~/.terraform.d/plugins/registry.terraform.io/ashleyjackson/instatus/0.1.0/darwin_arm64/
```

## Usage Example

### 1. Set Environment Variables

```bash
export INSTATUS_API_KEY="your-api-key"
```

### 2. Create Terraform Configuration

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
  # Uses INSTATUS_API_KEY environment variables
}

# Simple component
resource "instatus_component" "api" {
  name        = "API Service"
  description = "Our main API"
  status      = "OPERATIONAL"
  show_uptime = true
  order       = 1
}

# Parent group
resource "instatus_component" "website" {
  name        = "Website"
  description = "Main site"
  status      = "OPERATIONAL"
  show_uptime = true
  order       = 2
}

# Child component
resource "instatus_component" "promo" {
  name        = "Promo"
  description = "Promo site satellite"
  status      = "OPERATIONAL"
  show_uptime = true
  order       = 1
  grouped     = true
  group_id    = instatus_component.website.id
}
```

### 3. Apply Configuration

```bash
terraform init
terraform plan
terraform apply
```

## Testing Results

All CRUD operations tested successfully:

1. ✅ **Create** - Components created with correct attributes
2. ✅ **Read** - State correctly refreshed from API
3. ✅ **Update** - Component properties updated successfully
4. ✅ **Delete** - Components deleted cleanly
5. ✅ **Nested Hierarchy** - 3-level nesting working correctly:
   - Level 1: Test Services (parent)
   - Level 2: Test Databases (parent) + Test API (child)
   - Level 3: Primary DB + Replica DB (grandchildren)

## Component Attributes

### Required
- `name` - Component name

### Optional
- `description` - Component description
- `status` - Status: OPERATIONAL, UNDERMAINTENANCE, DEGRADEDPERFORMANCE, PARTIALOUTAGE, MAJOROUTAGE
- `show_uptime` - Show uptime metrics (default: true)
- `order` - Display order (default: 0)
- `grouped` - Whether component belongs to a group (default: false)
- `group_id` - Parent group ID (required if grouped=true)
- `archived` - Whether component is archived (default: false)

### Computed
- `id` - Component ID
- `unique_email` - Unique email for automation

## Next Steps

### Enhancements to Consider

1. **Data Sources** - Add data sources to query existing components
2. **Import** - Test importing existing components
3. **Validation** - Add more comprehensive input validation
4. **Status Updates** - Add resource for creating status incidents
5. **Bulk Operations** - Add support for bulk component management
6. **Testing** - Add unit and acceptance tests
7. **Publishing** - Publish to Terraform Registry

### Example: Adding a Data Source

Create `instatus/data_source_component.go`:

```go
func dataSourceComponent() *schema.Resource {
  return &schema.Resource{
    ReadContext: dataSourceComponentRead,
    Schema: map[string]*schema.Schema{
      "id": {
        Type:     schema.TypeString,
        Required: true,
      },
      // ... other attributes
    },
  }
}
```

## Troubleshooting

### Provider Not Found
- Ensure the provider is built and installed in the correct directory
- Check the architecture matches your system (darwin_arm64, darwin_amd64, linux_amd64)
- Run `terraform init` after installing

### API Errors
- Verify `INSTATUS_API_KEY` is set correctly
- Check API rate limits
- Ensure you have proper permissions

### Lock File Issues
```bash
rm -rf .terraform .terraform.lock.hcl
terraform init
```

## Documentation

Full documentation is available in the [README.md](README.md) file.

## License

MIT License
