# Web subdomain (app.kasoku.dev) → Railway
resource "google_dns_record_set" "web" {
  count        = var.enable_web_dns ? 1 : 0
  name         = "app.${var.domain}."
  managed_zone = google_dns_managed_zone.main.name
  project      = var.project_id
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["${var.web_dns_value}."]
}

# Railway domain verification for app.kasoku.dev
resource "google_dns_record_set" "web_verify" {
  count        = var.enable_web_dns && var.web_verify_txt != "" ? 1 : 0
  name         = "_railway-verify.app.${var.domain}."
  managed_zone = google_dns_managed_zone.main.name
  project      = var.project_id
  type         = "TXT"
  ttl          = 300
  rrdatas      = ["\"${var.web_verify_txt}\""]
}
