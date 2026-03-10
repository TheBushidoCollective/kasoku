# Server subdomain (api.kasoku.dev) → Railway
resource "google_dns_record_set" "server" {
  count        = var.enable_server_dns ? 1 : 0
  name         = "api.${var.domain}."
  managed_zone = google_dns_managed_zone.main.name
  project      = var.project_id
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["${var.server_dns_value}."]
}

# Railway domain verification for api.kasoku.dev
resource "google_dns_record_set" "server_verify" {
  count        = var.enable_server_dns && var.server_verify_txt != "" ? 1 : 0
  name         = "_railway-verify.api.${var.domain}."
  managed_zone = google_dns_managed_zone.main.name
  project      = var.project_id
  type         = "TXT"
  ttl          = 300
  rrdatas      = ["\"${var.server_verify_txt}\""]
}
