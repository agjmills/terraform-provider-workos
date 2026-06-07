resource "workos_organization" "example" {
  name = "Acme Corp"
  external_id = "a1b2c3d4"
  allow_profiles_outside_organization = false
  metadata = {
    tier = "enterprise"
  }
  domains {
    domain = "example.com"
    state  = "pending"
  }
}
