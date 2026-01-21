---
page_title: "instatus_page Resource - terraform-provider-instatus"
subcategory: ""
description: |-
  Manages an Instatus status page workspace.
---

# instatus_page (Resource)

Manages an Instatus status page workspace. Each page represents a separate status page with its own subdomain.

## Example Usage

### Basic Status Page

```terraform
resource "instatus_page" "example" {
  email          = "admin@example.com"
  name           = "Example Status Page"
  workspace_slug = "example-status"
}
```

### Status Page with Custom Branding

```terraform
resource "instatus_page" "branded" {
  email          = "admin@example.com"
  name           = "My Company Status"
  workspace_slug = "mycompany-status"
  
  logoUrl         = "https://example.com/logo.png"
  faviconUrl      = "https://example.com/favicon.ico"
  customDomain    = "status.mycompany.com"
  googleAnalytics = "UA-123456789-1"
}
```

### Multiple Pages

```terraform
variable "programs" {
  type    = list(string)
  default = ["app1", "app2", "app3"]
}

resource "instatus_page" "pages" {
  for_each = toset(var.programs)

  email          = "admin@example.com"
  name           = "Status Page for ${each.value}"
  workspace_slug = "${each.value}-status"
}
```

## Argument Reference

The following arguments are supported:

### Required

- `email` (String) - Billing email address for the status page workspace.
- `name` (String) - Display name for the status page workspace.
- `workspace_slug` (String) - Subdomain/slug for the status page. This will be used as `{workspace_slug}.instatus.com`. Cannot be changed after creation (forces new resource).

### Optional

- `logoUrl` (String) - URL of the logo to display on the status page.
- `faviconUrl` (String) - URL of the favicon for the status page.
- `googleAnalytics` (String) - Google Analytics tracking ID (e.g., `UA-XXXXXXXXX-X` or `G-XXXXXXXXXX`).
- `customDomain` (String) - Custom domain for the status page (e.g., `status.example.com`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) - The unique identifier for the status page.
- `workspace_id` (String) - The workspace ID returned by the Instatus API.

## Import

Status pages can be imported using the page ID:

```shell
terraform import instatus_page.example cmklxphmi0auy573fd29w4xoe
```

## Notes

- The `workspace_slug` becomes the subdomain: `{workspace_slug}.instatus.com`
- Custom domains require DNS configuration on your end
- Logo and favicon URLs must be publicly accessible
- The workspace is automatically created along with the page
- Deleting the page also deletes the associated workspace
