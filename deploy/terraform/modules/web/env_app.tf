resource "railway_variable" "node_env" {
  environment_id = var.environment_id
  service_id     = railway_service.web.id
  name           = "NODE_ENV"
  value          = "production"
}

resource "railway_variable" "api_url" {
  environment_id = var.environment_id
  service_id     = railway_service.web.id
  name           = "NEXT_PUBLIC_API_URL"
  value          = var.api_url
}
