output "service_id" {
  description = "Server service ID"
  value       = railway_service.server.id
}

output "service_name" {
  description = "Server service name"
  value       = railway_service.server.name
}

output "custom_domain_dns_value" {
  description = "DNS record value for custom domain (if configured)"
  value       = var.custom_domain != "" ? railway_custom_domain.server[0].dns_record_value : null
}
