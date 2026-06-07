# AGENTS.md

Guiding principles for agents contributing to this Terraform provider.

## Project Overview

This is a Terraform provider for [WorkOS](https://workos.com), enabling infrastructure-as-code management of WorkOS resources. Built with the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).

## Technology Stack

- **Language:** Go 1.21+
- **Framework:** `terraform-plugin-framework` (not SDKv2)
- **API Client:** `github.com/workos/workos-go/v3` (v4+ when available)
- **Testing:** `github.com/hashicorp/terraform-plugin-testing` + standard Go tests
- **Linting:** `golangci-lint`
- **Release:** `goreleaser` with conventional commits determining semver bumps

## Project Conventions

### File Organization
```
internal/provider/
├── provider.go                        # Provider definition, schema, Configure
├── resource_organization.go           # workos_organization resource
├── resource_webhook_endpoint.go       # workos_webhook_endpoint resource
├── resource_redirect_uri.go           # workos_redirect_uri resource
├── resource_cors_origin.go            # workos_cors_origin resource
├── data_source_organization.go        # data "workos_organization"
├── data_source_sso_connection.go      # data "workos_sso_connection"
├── data_source_directory.go           # data "workos_directory"
├── data_source_webhook_endpoint.go    # data "workos_webhook_endpoint"
├── data_source_organizations.go       # plural list data sources
├── data_source_sso_connections.go
├── data_source_directories.go
├── data_source_webhook_endpoints.go
├── *__test.go                         # Test files co-located with source
├── test_utils.go                      # Shared test helpers
└── model_*.go                         # Shared type models (if needed)
```

### Resource / Data Source Patterns

Every resource and data source follows this pattern:
1. Define the `*Resource` / `*DataSource` struct implementing the appropriate interface
2. Schema definition in the `Schema` method using the framework's `schema.Schema` type
3. Model struct using `tfsdk.Attribute` / `schema.Attribute` tags
4. Validators for required fields
5. Null/known value checking in CRUD methods

### Naming Conventions

- Terraform resources: `workos_<resource_name>` in snake_case
- Data sources: `workos_<resource_name>` (singular) + `workos_<resource_names>` (plural list)
- Go types: CamelCase matching WorkOS API field names
- Computed-only fields: `id`, `created_at`, `updated_at`, `secret`
- Sensitive fields: `secret`, any key/token fields

### API Client

- Use `github.com/workos/workos-go/v3` SDK
- Wrap in `internal/client/client.go` for provider-wide configuration
- Client is configured with API key at provider init, stored in provider struct
- All API calls go through the client wrapper

### Error Handling

- Always check errors from WorkOS API calls
- Use `diag.Diagnostics` for Terraform-facing errors
- Add descriptive error messages: `"Unable to create WorkOS organization: " + err.Error()`
- 404 on read = remove from state (set to null)

### Commit Conventions

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` - new resource, data source, or significant functionality
- `fix:` - bug fix
- `docs:` - documentation changes
- `chore:` - build, CI, dependencies
- `test:` - test additions/changes
- `refactor:` - code restructuring without functionality change

Breaking changes: add `!` after the type (`feat!:` / `fix!:`) or add `BREAKING CHANGE:` footer.

### Release Process

1. Commits to `main` follow conventional commits
2. `goreleaser` runs on tag pushes
3. Semantic version is derived from conventional commit messages
4. Release artifacts: binaries for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64
5. GPG signing via GitHub Actions

## Testing Guidelines

### Unit Tests
- Mock the WorkOS client interface
- Test schema validation
- Test model conversion functions

### Acceptance Tests
- TF_ACC=1 environment variable
- Use `WORKOS_API_KEY` from environment
- Tests create real resources in a WorkOS sandbox
- Clean up resources in `CheckDestroy` functions
- Use `t.Parallel()` when tests are independent

### Test Patterns
```go
func TestAccResourceOrganization_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccResourceOrganizationConfig("Test Org"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("workos_organization.test", "name", "Test Org"),
                    resource.TestCheckResourceAttrSet("workos_organization.test", "id"),
                ),
            },
        },
    })
}
```

## Documentation

- Generate docs using `tfplugindocs` or maintain manually in `docs/`
- Each resource/data source needs an index entry in `docs/index.md`
- Examples in `examples/` directory
- Schema is source of truth for attribute documentation

## Resources

| Resource | WorkOS API | CRUD Support |
|---|---|---|
| `workos_organization` | `/organizations` | Full CRUD |
| `workos_webhook_endpoint` | `/webhook_endpoints` | Full CRUD |
| `workos_redirect_uri` | `/user_management/redirect_uris` | Create, Read, Delete |
| `workos_cors_origin` | `/user_management/cors_origins` | Create, Read, Delete |

## Data Sources

| Data Source | WorkOS API | Lookup By |
|---|---|---|
| `workos_organization` | GET `/organizations/{id}` or `/external_id/{id}` | id, external_id |
| `workos_organizations` | GET `/organizations` | domains, search, limit |
| `workos_sso_connection` | GET `/connections/{id}` | id |
| `workos_sso_connections` | GET `/connections` | org_id, type, domain, search |
| `workos_directory` | GET `/directories/{id}` | id |
| `workos_directories` | GET `/directories` | org_id, domain, search |
| `workos_webhook_endpoint` | GET `/webhook_endpoints/{id}` | id |
| `workos_webhook_endpoints` | GET `/webhook_endpoints` | -- |
