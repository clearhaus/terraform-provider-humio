resource "humio_action" "example_opsgenie" {
  repository = "sandbox"
  name       = "example_opsgenie"
  type       = "OpsGenieAction"

  opsgenie {
    api_url   = "https://api.opsgenie.com"
    genie_key = "XXXXXXXXXXXXXXX"
  }
}
