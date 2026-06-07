---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_webhook_endpoint"
description: |-
  Manages a WorkOS webhook endpoint.
---

# Resource: workos_webhook_endpoint

Manages a WorkOS webhook endpoint. Webhook endpoints receive event notifications from WorkOS.

## Example Usage

```terraform
resource "workos_webhook_endpoint" "example" {
  endpoint_url = "https://example.com/webhooks"
  events       = ["user.created", "dsync.user.created"]
  status       = "enabled"
}
```

## Argument Reference

- `endpoint_url` (Required, String) - The URL that will receive webhook events.
- `events` (Required, List of String) - The event types this endpoint subscribes to.
- `status` (Optional, String) - The status: `enabled` or `disabled`. Defaults to `enabled`.

## Attribute Reference

- `id` (String) - The unique identifier.
- `secret` (String, Sensitive) - The signing secret for webhook payloads.
- `created_at` (String) - ISO 8601 timestamp of creation.
- `updated_at` (String) - ISO 8601 timestamp of last update.

## Import

```shell
terraform import workos_webhook_endpoint.example we_01HXYZ123456789ABCDEFGHIJ
```
