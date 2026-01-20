---
page_title: "instatus_page Resource - terraform-provider-instatus"
subcategory: ""
description: |-
  Manages an Instatus status page workspace.
---

# instatus_page (Resource)

Manages an Instatus status page workspace. Each page represents a separate status page with its own subdomain.

## Example Usage

```terraform
resource "instatus_page" "example" {
  email          = "admin@example.com"
  name           = "Example Status Page"
  workspace_slug = "example-status"
}
```

## Multiple Pages

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

- `email` - (Required) Billing email address for the status page workspace.
- `name` - (Required) Display name for the status page workspace.
- `workspace_slug` - (Required) Subdomain/slug for the status page. This will be used as `{workspace_slug}.instatus.com`. Cannot be changed after creation (forces new resource).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the status page.
- `workspace_id` - The workspace ID returned by the Instatus API.

## Import

Status pages can be imported using the page ID:

```shell
terraform import instatus_page.example cmklxphmi0auy573fd29w4xoe
```
