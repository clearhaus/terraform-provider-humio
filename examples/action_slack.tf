resource "humio_action" "example_slack" {
  repository = "sandbox"
  name       = "example_slack"
  type       = "SlackAction"

  slack {
    url = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    fields = {
      "Events String" = "{events_str}"
      "Query"         = "{query_string}"
      "Time Interval" = "{query_time_interval}"
    }
  }
}
