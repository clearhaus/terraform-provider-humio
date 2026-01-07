resource "humio_alert" "example_alert_with_labels" {
  repository = humio_action.example_email.repository
  name       = "example_alert_with_labels"

  actions = [humio_action.example_email.action_id]

  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "12h"
}

resource "humio_alert" "example_alert_with_throttle_field" {
  repository = humio_action.example_email.repository
  name       = "example_alert_with_throttle_field"

  actions = [humio_action.example_email.action_id]

  throttle_time_millis = 300000
  throttle_field       = "serviceName"
  enabled              = true
  query                = "level = ERROR"
  start                = "12h"
}

resource "humio_alert" "example_alert_without_labels" {
  repository = humio_action.example_email_body.repository
  name       = "example_alert_without_labels"

  actions = [humio_action.example_email_body.action_id]

  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "12h"
}

resource "humio_alert" "example_alert_with_description" {
  repository  = humio_action.example_email_body.repository
  name        = "example_alert_with_description"
  description = "lorem ipsum...."

  actions = [
    humio_action.example_email_body.action_id,
    humio_action.example_email_subject.action_id,
  ]

  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "1d"
}

resource "humio_alert" "example_alert_with_user_owner" {
  repository = humio_action.example_email_body.repository
  name       = "example_alert_with_user_owner"

  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "1d"

  # run_as_user_id       = "XXXXXXXXXXXXXXXXXXXXXXXX"
  query_ownership_type = "User"
}

resource "humio_repository" "org_ownership_test" {
  name        = "example_repo_org_ownership_${local.email_prefix}"
  description = "Repository for organization query ownership"
}

resource "humio_alert" "example_alert_with_organization_owner" {
  repository = humio_repository.org_ownership_test.name
  name       = "example_alert_with_organization_owner"

  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "1d"

  query_ownership_type = "Organization"
}
