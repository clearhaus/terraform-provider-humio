resource "humio_action" "example_pagerduty" {
  repository = "sandbox"
  name       = "example_pagerduty"
  type       = "PagerDutyAction"

  pagerduty {
    routing_key = "XXXXXXXXXXXXXXX"
    severity    = "critical"
  }
}
