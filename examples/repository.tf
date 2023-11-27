resource "humio_repository" "example_repo_minimal_fields_set" {
  name = "example_repo_minimal_fields_set"
}

resource "humio_repository" "example_repo_all_fields_set" {
  name        = "example_repo_all_fields_set"
  description = "This is an example"

  retention {
    time_in_days = 30
  }
}
