
resource "google_compute_network" "network" {
    project = "${var.project}"
    name = "${var.name}"
    auto_create_subnetworks = false
}

variable name {
  type        = string
  description = "network name"
}

variable project {
  type        = string
  description = "project name"
}
