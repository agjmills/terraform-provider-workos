---
subcategory: ""
layout: "workos"
page_title: "WorkOS: workos_webhook_endpoints"
description: |-
  List WorkOS webhook endpoints.
---

# Data Source: workos_webhook_endpoints

List all WorkOS webhook endpoints.

## Example Usage

```terraform
data "workos_webhook_endpoints" "example" {}
```

## Attribute Reference

- `webhook_endpoints` (List of Object) - The list of webhook endpoints.
- `list_metadata` (Object) - Pagination cursors.
