---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_organization"
description: |-
  Manages a WorkOS organization.
---

# Resource: workos_organization

Manages a WorkOS organization. Organizations are top-level resources in WorkOS. Each Connection, Directory, and Audit Trail Event belongs to an Organization.

## Example Usage

```terraform
resource "workos_organization" "example" {
  name = "Acme Corp"
  external_id = "a1b2c3d4"
  allow_profiles_outside_organization = false
  metadata = {
    tier = "enterprise"
  }
  domains {
    domain = "acme.com"
    state  = "pending"
  }
}
```

## Argument Reference

- `name` (Required, String) - The name of the organization.
- `external_id` (Optional, String) - An external identifier for the organization.
- `allow_profiles_outside_organization` (Optional, Bool) - Whether profiles outside the organization are allowed.
- `metadata` (Optional, Map of String) - Key-value metadata for the organization.
- `domains` (Optional, List of Object) - Organization domains. Each object has:
  - `domain` (Required, String) - Domain name.
  - `state` (Optional, String) - Domain verification state (e.g., `pending`, `verified`).
  - `verification_strategy` (Optional, String) - Domain verification strategy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) - The unique identifier of the organization.
- `created_at` (String) - ISO 8601 timestamp of creation.
- `updated_at` (String) - ISO 8601 timestamp of last update.

## Import

Organizations can be imported by ID:

```shell
terraform import workos_organization.example org_01HXYZ123456789ABCDEFGHIJ
```
