# Get information about the current authenticated user
data "humio_user" "current" {}

# Example: Use the email prefix (before @) in a repository name
locals {
  email_prefix = lower(replace(split("@", data.humio_user.current.email)[0], "/[^a-z0-9]/", ""))
}

# Output user information
output "current_user_id" {
  value       = data.humio_user.current.id
  description = "The ID of the current authenticated user"
}

output "current_username" {
  value       = data.humio_user.current.username
  description = "The username of the current authenticated user"
}

output "current_user_email" {
  value       = data.humio_user.current.email
  description = "The email of the current authenticated user"
}

output "current_user_full_name" {
  value       = data.humio_user.current.full_name
  description = "The full name of the current authenticated user"
}

output "is_root_user" {
  value       = data.humio_user.current.is_root
  description = "Whether the current user is a root user"
}

