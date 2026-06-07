---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_directories"
description: |-
  List WorkOS directories.
---

# Data Source: workos_directories

List WorkOS directories with optional filtering.

## Example Usage

```terraform
data "workos_directories" "example" {
  organization_id = "org_01HXYZ123456789ABCDEFGHIJ"
  search          = "corp"
}
```

## Argument Reference

- `organization_id` (Optional, String) - Filter by organization.
- `domain` (Optional, String) - Filter by domain.
- `search` (Optional, String) - Search text.

## Attribute Reference

- `directories` (List of Object) - The list of directories.
- `list_metadata` (Object) - Pagination cursors.
