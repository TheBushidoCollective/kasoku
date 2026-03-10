resource "random_password" "jwt_secret" {
  length  = 64
  special = false
}

resource "railway_variable" "jwt_secret" {
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "JWT_SECRET"
  value          = random_password.jwt_secret.result
}

resource "railway_variable" "github_client_id" {
  count          = var.github_client_id != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "GITHUB_CLIENT_ID"
  value          = var.github_client_id
}

resource "railway_variable" "github_client_secret" {
  count          = var.github_client_secret != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "GITHUB_CLIENT_SECRET"
  value          = var.github_client_secret
}

resource "railway_variable" "gitlab_client_id" {
  count          = var.gitlab_client_id != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "GITLAB_CLIENT_ID"
  value          = var.gitlab_client_id
}

resource "railway_variable" "gitlab_client_secret" {
  count          = var.gitlab_client_secret != "" ? 1 : 0
  environment_id = var.environment_id
  service_id     = railway_service.server.id
  name           = "GITLAB_CLIENT_SECRET"
  value          = var.gitlab_client_secret
}
