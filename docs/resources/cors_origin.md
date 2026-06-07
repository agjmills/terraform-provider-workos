---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_cors_origin"
description: |-
  Manages a WorkOS AuthKit CORS origin.
---

# Resource: workos_cors_origin

Manages a WorkOS AuthKit CORS origin.

## Example Usage

```terraform
resource "workos_cors_origin" "example" {
  origin = "https://example.com"
}
```

## Argument Reference

- `origin` (Required, String) - The CORS origin URL.

## Attribute Reference

- `id` (String) - The unique identifier.
- `created_at` (String) - ISO 8601 timestamp of creation.
- `updated_at` (String) - ISO 8601 timestamp of last update.

## Import

```shell
terraform import workos_cors_origin.example cors_origin_01HXYZ123456789ABCDEFGHIJ
```
