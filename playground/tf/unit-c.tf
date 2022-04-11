// =================================
// Terraform Module Config: Unit "C"
// =================================
//
// This file pickups configurations from the "project" variable
// to plan resources for "unit-c" for the configfured stages.
//

locals {
  unit_c_projects = {
    for proj, conf in local.projects_flattened: proj => conf if conf.name == "c"
  }
}

module "unit-c" {
  source = "./net"
  for_each = local.unit_c_projects
  name = each.value.name
  project = each.value.project
}
