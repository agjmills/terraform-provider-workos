---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_organizations"
description: |-
  List WorkOS organizations.
---

# Data Source: workos_organizations

List WorkOS organizations with optional filtering.

## Example Usage

```terraform
data "workos_organizations" "example" {
  domains = ["acme.com"]
  search  = "Acme"
}
```

## Argument Reference

- `domains` (Optional, List of String) - Filter by associated domains.
- `search` (Optional, String) - Search text to match against organization names.

## Attribute Reference

- `organizations` (List of Object) - The list of organizations.
- `list_metadata` (Object) - Pagination cursors (`before`, `after`).
