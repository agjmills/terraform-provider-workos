---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_redirect_uri"
description: |-
  Manages a WorkOS AuthKit redirect URI.
---

# Resource: workos_redirect_uri

Manages a WorkOS AuthKit redirect URI.

## Example Usage

```terraform
resource "workos_redirect_uri" "example" {
  uri = "https://example.com/callback"
}
```

## Argument Reference

- `uri` (Required, String) - The redirect URI.

## Attribute Reference

- `id` (String) - The unique identifier.
- `default` (Bool) - Whether this is the default redirect URI.
- `created_at` (String) - ISO 8601 timestamp of creation.
- `updated_at` (String) - ISO 8601 timestamp of last update.

## Import

```shell
terraform import workos_redirect_uri.example ruri_01HXYZ123456789ABCDEFGHIJ
```
