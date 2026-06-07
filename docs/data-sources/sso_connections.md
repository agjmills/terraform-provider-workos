---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_sso_connections"
description: |-
  List WorkOS SSO connections.
---

# Data Source: workos_sso_connections

List WorkOS SSO connections with optional filtering.

## Example Usage

```terraform
data "workos_sso_connections" "example" {
  organization_id = "org_01HXYZ123456789ABCDEFGHIJ"
  connection_type = "OktaSAML"
}
```

## Argument Reference

- `organization_id` (Optional, String) - Filter by organization.
- `connection_type` (Optional, String) - Filter by connection type.
- `domain` (Optional, String) - Filter by domain.
- `search` (Optional, String) - Search text.

## Attribute Reference

- `connections` (List of Object) - The list of connections.
- `list_metadata` (Object) - Pagination cursors.
