---
layout: ""
page_title: "WorkOS Provider"
description: |-
  The WorkOS provider provides resources to interact with the WorkOS API.
---

# WorkOS Provider

The WorkOS provider is used to configure and manage resources in [WorkOS](https://workos.com).

## Example Usage

```terraform
terraform {
  required_providers {
    workos = {
      source = "workos/workos"
    }
  }
}

provider "workos" {
  api_key = "sk_example_123456789"
}
```

## Authentication

The provider requires a WorkOS API key. You can provide it via:

1. **Environment variable**: Set `WORKOS_API_KEY` to your API key.
2. **Provider configuration**: Use the `api_key` attribute in the provider block.

## Schema

### Optional

- `api_key` (String, Sensitive) - The WorkOS API key. Can also be set via the `WORKOS_API_KEY` environment variable.
