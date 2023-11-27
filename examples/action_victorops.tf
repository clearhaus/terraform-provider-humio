resource "humio_action" "example_victorops" {
  repository = "sandbox"
  name       = "example_victorops"
  type       = "VictorOpsAction"

  victorops {
    message_type = "critical"
    notify_url   = "https://example.org"
  }
}
