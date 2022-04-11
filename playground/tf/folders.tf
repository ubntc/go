locals {
  // ================
  // Helper Variables
  // ================
  //
  // The following lists and maps are needed to iterate over the nested
  // structures in the "projects" variable and to ensure only known values
  // are used.

  // lookup table to ensure only valid stages are configured
  stages_lookup = {
    live     = "live"
    unstable = "unstable"
    staging  = "staging"
    sandbox  = "sandbox"
  }

  // lookup table to ensure only valid project names are configured
  project_lookup = {
    unit-a = "a"
    unit-b = "b"
    unit-c = "c"
  }

  projects_flattened = merge([
    for parent, items in var.projects: { for name in keys(items): "${parent}.${name}" => merge(
      {
        parent = local.stages_lookup[parent],
        name = local.project_lookup[name],
        project = "${local.project_lookup[name]}-${local.stages_lookup[parent]}",
      },
      var.projects[parent][name],
    )}
  ]...)
}

module "stage_folders" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 3.0"

  parent = var.organization_id
  names  = [for k in keys(var.projects) : local.stages_lookup[k]]

  set_roles = false
}
