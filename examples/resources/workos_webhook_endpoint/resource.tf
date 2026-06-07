resource "workos_webhook_endpoint" "example" {
  endpoint_url = "https://example.com/webhooks"
  events       = ["user.created", "dsync.user.created"]
}
