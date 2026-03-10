variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "zone_name" {
  description = "Name of the Cloud DNS managed zone"
  type        = string
}

variable "domain" {
  description = "Domain name (e.g., kasoku.dev)"
  type        = string
}

# Server subdomain
variable "enable_server_dns" {
  description = "Whether to create server subdomain DNS records"
  type        = bool
  default     = false
}

variable "server_dns_value" {
  description = "Railway DNS value for server subdomain"
  type        = string
  default     = ""
}

variable "server_verify_txt" {
  description = "Railway domain verification TXT value for server"
  type        = string
  default     = ""
}

# Web subdomain
variable "enable_web_dns" {
  description = "Whether to create web subdomain DNS records"
  type        = bool
  default     = false
}

variable "web_dns_value" {
  description = "Railway DNS value for web subdomain"
  type        = string
  default     = ""
}

variable "web_verify_txt" {
  description = "Railway domain verification TXT value for web"
  type        = string
  default     = ""
}
