resource "railway_variable" "stripe_secret_key" {
  count          = var.stripe_secret_key != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "STRIPE_SECRET_KEY"
  value          = var.stripe_secret_key
}

resource "railway_variable" "stripe_webhook_secret" {
  count          = var.stripe_webhook_secret != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "STRIPE_WEBHOOK_SECRET"
  value          = var.stripe_webhook_secret
}
