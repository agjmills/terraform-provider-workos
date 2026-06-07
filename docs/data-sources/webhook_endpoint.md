---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_webhook_endpoint"
description: |-
  Get information about a WorkOS webhook endpoint.
---

# Data Source: workos_webhook_endpoint

Use this data source to get information about a WorkOS webhook endpoint.

## Example Usage

```terraform
data "workos_webhook_endpoint" "example" {
  id = "we_01HXYZ123456789ABCDEFGHIJ"
}
```

## Argument Reference

- `id` (Required, String) - The webhook endpoint ID.

## Attribute Reference

- `endpoint_url` (String) - The webhook URL.
- `secret` (String, Sensitive) - Signing secret.
- `status` (String) - Status (`enabled` or `disabled`).
- `events` (List of String) - Subscribed events.
- `created_at` (String)
- `updated_at` (String)
