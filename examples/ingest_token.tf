resource "humio_ingest_token" "example_ingest_token_without_parser" {
  repository = "sandbox"
  name       = "example_ingest_token_without_parser"
}

resource "humio_ingest_token" "example_ingest_token_with_accesslog_parser" {
  repository = "sandbox"
  name       = "example_ingest_token_with_accesslog_parser"
  parser     = "accesslog"

  depends_on = [
    humio_ingest_token.example_ingest_token_without_parser
  ]
}

output "ingest_token_without_parser" {
  value     = humio_ingest_token.example_ingest_token_without_parser.token
  sensitive = true
}

output "ingest_token_with_accesslog_parser" {
  value     = humio_ingest_token.example_ingest_token_with_accesslog_parser.token
  sensitive = true
}
