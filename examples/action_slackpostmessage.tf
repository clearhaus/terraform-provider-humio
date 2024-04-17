resource "humio_action" "example_slackpostmessage" {
  repository = "sandbox"
  name       = "example_slackpostmessage"
  type       = "SlackPostMessageAction"

  slackpostmessage {
    api_token = "abcdefghij1234567890"
    channels  = ["#alerts", "#ops"]
    fields = {
      "Events String" = "{events_str}"
      "Query"         = "{query_string}"
      "Time Interval" = "{query_time_interval}"
    }
  }
}
