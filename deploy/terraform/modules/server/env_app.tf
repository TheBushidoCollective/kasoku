resource "railway_variable" "port" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "PORT"
  value          = "8080"
}

resource "railway_variable" "storage_type" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "STORAGE_TYPE"
  value          = "local"
}

resource "railway_variable" "storage_path" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "STORAGE_PATH"
  value          = "/app/storage"
}

resource "railway_variable" "base_url" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "BASE_URL"
  value          = var.custom_domain != "" ? "https://${var.custom_domain}" : ""
}
