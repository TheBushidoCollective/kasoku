resource "railway_variable" "db_driver" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "DB_DRIVER"
  value          = "postgres"
}

resource "railway_variable" "db_dsn" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "DB_DSN"
  value          = var.database_url
}
