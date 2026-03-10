# -----------------------------------------------------------------------------
# Outputs
# -----------------------------------------------------------------------------

output "project_id" {
  description = "Railway project ID"
  value       = module.project.id
}

output "server_service_id" {
  description = "Server service ID"
  value       = module.server.service_id
}

output "web_service_id" {
  description = "Web service ID"
  value       = module.web.service_id
}

output "postgres_service_id" {
  description = "PostgreSQL service ID"
  value       = module.postgres.service_id
}

output "server_url" {
  description = "Server URL (Railway-generated or custom domain)"
  value       = var.server_custom_domain != "" ? "https://${var.server_custom_domain}" : "https://${module.server.service_name}.up.railway.app"
}

output "web_url" {
  description = "Web URL (Railway-generated or custom domain)"
  value       = var.web_custom_domain != "" ? "https://${var.web_custom_domain}" : "https://${module.web.service_name}.up.railway.app"
}

output "pr_environments_enabled" {
  description = "Whether PR preview environments are enabled"
  value       = var.enable_pr_environments
}
