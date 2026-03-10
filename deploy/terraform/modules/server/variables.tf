# Required
variable "project_id" {
  description = "Railway project ID"
  type        = string
}

variable "environment_id" {
  description = "Railway environment ID"
  type        = string
}

variable "database_url" {
  description = "PostgreSQL connection URL"
  type        = string
  sensitive   = true
}

# Service Configuration
variable "service_name" {
  description = "Name of the server service"
  type        = string
  default     = "server"
}

variable "github_repo" {
  description = "GitHub repository (owner/repo format)"
  type        = string
}

variable "branch" {
  description = "Git branch to deploy"
  type        = string
  default     = "main"
}

variable "root_directory" {
  description = "Root directory for the service"
  type        = string
  default     = "server"
}

# OAuth (optional)
variable "github_client_id" {
  description = "GitHub OAuth client ID"
  type        = string
  sensitive   = true
  default     = ""
}

variable "github_client_secret" {
  description = "GitHub OAuth client secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "gitlab_client_id" {
  description = "GitLab OAuth client ID"
  type        = string
  sensitive   = true
  default     = ""
}

variable "gitlab_client_secret" {
  description = "GitLab OAuth client secret"
  type        = string
  sensitive   = true
  default     = ""
}

# Billing (optional)
variable "stripe_secret_key" {
  description = "Stripe secret key (empty to disable billing)"
  type        = string
  sensitive   = true
  default     = ""
}

variable "stripe_webhook_secret" {
  description = "Stripe webhook secret"
  type        = string
  sensitive   = true
  default     = ""
}

# Custom Domain
variable "custom_domain" {
  description = "Custom domain for the server (optional)"
  type        = string
  default     = ""
}
