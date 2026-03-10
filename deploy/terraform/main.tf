# -----------------------------------------------------------------------------
# Kasoku - Railway Infrastructure
# -----------------------------------------------------------------------------

module "project" {
  source = "./modules/project"

  name           = var.project_name
  description    = "Kasoku - Build caching server and web dashboard"
  has_pr_deploys = var.enable_pr_environments
}

module "postgres" {
  source = "./modules/postgres"

  project_id     = module.project.id
  environment_id = module.project.default_environment_id
  database       = "kasoku"
  username       = "kasoku"
}

module "server" {
  source = "./modules/server"

  project_id     = module.project.id
  environment_id = module.project.default_environment_id

  # GitHub source
  github_repo    = var.github_repo
  branch         = var.production_branch
  root_directory = var.server_root_directory

  # Database
  database_url = module.postgres.connection_url

  # OAuth (optional)
  github_client_id     = var.github_client_id
  github_client_secret = var.github_client_secret
  gitlab_client_id     = var.gitlab_client_id
  gitlab_client_secret = var.gitlab_client_secret

  # Billing (optional)
  stripe_secret_key     = var.stripe_secret_key
  stripe_webhook_secret = var.stripe_webhook_secret

  # Custom domain
  custom_domain = var.server_custom_domain
}

module "web" {
  source = "./modules/web"

  project_id     = module.project.id
  environment_id = module.project.default_environment_id

  # GitHub source
  github_repo    = var.github_repo
  branch         = var.production_branch
  root_directory = var.web_root_directory

  # API URL (use custom domain if set, otherwise Railway internal DNS)
  api_url = var.server_custom_domain != "" ? "https://${var.server_custom_domain}" : "http://server.railway.internal:8080"

  # Custom domain
  custom_domain = var.web_custom_domain
}

# -----------------------------------------------------------------------------
# GCP DNS (optional - only if gcp_project_id is set)
# -----------------------------------------------------------------------------

module "dns" {
  source = "./modules/dns"
  count  = var.gcp_project_id != "" ? 1 : 0

  project_id = var.gcp_project_id
  zone_name  = var.gcp_dns_zone_name
  domain     = var.domain

  # Server subdomain (api.kasoku.dev) → Railway
  enable_server_dns = var.server_custom_domain != ""
  server_dns_value  = module.server.custom_domain_dns_value
  server_verify_txt = var.server_domain_verify_txt

  # Web subdomain (app.kasoku.dev) → Railway
  enable_web_dns = var.web_custom_domain != ""
  web_dns_value  = module.web.custom_domain_dns_value
  web_verify_txt = var.web_domain_verify_txt
}
