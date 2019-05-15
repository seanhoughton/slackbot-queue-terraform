module "relay" {
  source       = "../terraform"
  service_name = "${var.service_name}"
}
