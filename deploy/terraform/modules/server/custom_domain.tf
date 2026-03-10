resource "railway_custom_domain" "server" {
  count          = var.custom_domain != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  domain         = var.custom_domain
}
