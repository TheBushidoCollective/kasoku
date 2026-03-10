# -----------------------------------------------------------------------------
# Project Configuration
# -----------------------------------------------------------------------------

variable "project_name" {
  description = "Name of the Railway project"
  type        = string
  default     = "kasoku"
}

variable "github_repo" {
  description = "GitHub repository for deployments (owner/repo format)"
  type        = string
  default     = "TheBushidoCollective/kasoku"
}

variable "production_branch" {
  description = "Branch to deploy to production"
  type        = string
  default     = "main"
}

# -----------------------------------------------------------------------------
# Environment Configuration
# -----------------------------------------------------------------------------

variable "enable_pr_environments" {
  description = "Enable automatic PR preview environments"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# Server Service Configuration
# -----------------------------------------------------------------------------

variable "server_root_directory" {
  description = "Root directory for the server service"
  type        = string
  default     = "server"
}

variable "server_custom_domain" {
  description = "Custom domain for the server API (e.g., api.kasoku.dev)"
  type        = string
  default     = ""
}

# -----------------------------------------------------------------------------
# Web Service Configuration
# -----------------------------------------------------------------------------

variable "web_root_directory" {
  description = "Root directory for the web service"
  type        = string
  default     = "web"
}

variable "web_custom_domain" {
  description = "Custom domain for the web dashboard (e.g., app.kasoku.dev)"
  type        = string
  default     = ""
}

# -----------------------------------------------------------------------------
# GCP DNS Configuration (optional)
# -----------------------------------------------------------------------------

variable "gcp_project_id" {
  description = "GCP project ID for Cloud DNS (empty to skip DNS)"
  type        = string
  default     = ""
}

variable "gcp_dns_zone_name" {
  description = "Name of the Cloud DNS managed zone"
  type        = string
  default     = "kasoku-dev"
}

variable "domain" {
  description = "Base domain"
  type        = string
  default     = "kasoku.dev"
}

variable "server_domain_verify_txt" {
  description = "Railway domain verification TXT value for server subdomain"
  type        = string
  default     = ""
}

variable "web_domain_verify_txt" {
  description = "Railway domain verification TXT value for web subdomain"
  type        = string
  default     = ""
}

# -----------------------------------------------------------------------------
# OAuth Secrets (sensitive - pass via TF_VAR_* or terraform.tfvars)
# -----------------------------------------------------------------------------

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

# -----------------------------------------------------------------------------
# Stripe Billing (optional - pass via TF_VAR_* or terraform.tfvars)
# -----------------------------------------------------------------------------

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
