module "swarmbot" {
  source       = "../terraform"
  service_name = "${var.service_name}"
}
