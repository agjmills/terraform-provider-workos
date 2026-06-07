---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_directory"
description: |-
  Get information about a WorkOS directory.
---

# Data Source: workos_directory

Use this data source to get information about a WorkOS directory.

## Example Usage

```terraform
data "workos_directory" "example" {
  id = "directory_01HXYZ123456789ABCDEFGHIJ"
}
```

## Argument Reference

- `id` (Required, String) - The directory ID.

## Attribute Reference

- `organization_id` (String) - Parent organization ID.
- `external_key` (String) - External key.
- `type` (String) - Directory type.
- `state` (String) - Directory state.
- `name` (String) - Directory name.
- `domain` (String) - Associated domain.
- `metadata` (Map of String) - Aggregate counts.
- `created_at` (String)
- `updated_at` (String)
