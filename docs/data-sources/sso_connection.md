---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_sso_connection"
description: |-
  Get information about a WorkOS SSO connection.
---

# Data Source: workos_sso_connection

Use this data source to get information about a WorkOS SSO connection.

## Example Usage

```terraform
data "workos_sso_connection" "example" {
  id = "conn_01HXYZ123456789ABCDEFGHIJ"
}
```

## Argument Reference

- `id` (Required, String) - The connection ID.

## Attribute Reference

- `organization_id` (String) - Parent organization ID.
- `connection_type` (String) - Connection type (e.g., `OktaSAML`).
- `name` (String) - Connection name.
- `state` (String) - Connection state (`active`, `inactive`, `validating`).
- `domains` (List of Object) - Associated domains.
- `options` (Map of String) - Connection options.
- `created_at` (String)
- `updated_at` (String)
