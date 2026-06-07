# Terraform Provider WorkOS

Terraform provider for [WorkOS](https://workos.com) - manage your WorkOS resources as infrastructure as code.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (to build the provider plugin)

## Using the Provider

```hcl
terraform {
  required_providers {
    workos = {
      source  = "workos/workos"
      version = "~> 1.0"
    }
  }
}

provider "workos" {
  api_key = var.workos_api_key
}
```

## Resources

| Resource | Description |
|---|---|
| `workos_organization` | Manage WorkOS organizations |
| `workos_webhook_endpoint` | Manage WorkOS webhook endpoints |
| `workos_redirect_uri` | Manage AuthKit redirect URIs |
| `workos_cors_origin` | Manage AuthKit CORS origins |

## Data Sources

| Data Source | Description |
|---|---|
| `workos_organization` | Look up an organization by ID or external ID |
| `workos_organizations` | List organizations with filtering |
| `workos_sso_connection` | Look up an SSO connection by ID |
| `workos_sso_connections` | List SSO connections with filtering |
| `workos_directory` | Look up a directory by ID |
| `workos_directories` | List directories with filtering |
| `workos_webhook_endpoint` | Look up a webhook endpoint by ID |
| `workos_webhook_endpoints` | List webhook endpoints |

## Development

```bash
make build    # Build the provider
make test     # Run unit tests
make testacc  # Run acceptance tests (requires WORKOS_API_KEY)
make vet      # Run go vet
```

## License

MIT
