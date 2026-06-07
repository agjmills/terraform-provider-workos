---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_organization"
description: |-
  Get information about a WorkOS organization.
---

# Data Source: workos_organization

Use this data source to get information about a WorkOS organization by ID or external ID.

## Example Usage

```terraform
data "workos_organization" "example" {
  id = "org_01HXYZ123456789ABCDEFGHIJ"
}
```

## Argument Reference

- `id` (Optional, String) - The organization ID.
- `external_id` (Optional, String) - The external identifier. Exactly one of `id` or `external_id` must be specified.

## Attribute Reference

- `name` (String) - The organization name.
- `domains` (List of Object) - Organization domains.
- `metadata` (Map of String) - Key-value metadata.
- `allow_profiles_outside_organization` (Bool)
- `stripe_customer_id` (String)
- `created_at` (String)
- `updated_at` (String)
