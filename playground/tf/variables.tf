// ==============
// Main Variables
// ==============
//
// Variables define configuration option for our platform.
// Not all resources are reflected in the variables.
// Some resources may by automatically defined by different modules for each
// configured "stage" or "project".

variable projects {
  type = map(map(object({
    cidr_block = string
  })))

  description = "folders|environments|stages > projects|units|components > values"
}

variable organization_id {
  type        = string
  description = "Google Organization ID"
}

// ======================
// Local Helper Variables
// ======================
//
// The following lists and maps are needed to iterate over the nested
// structures in the "projects" variable and to ensure only known values
// are used.

locals {
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

  // flattened variant of the "projects" variable.
  //   input:   { "stage": { "unit": values }, ... }
  //   output:  { "stage.unit": values, ... }
  //
  // The derived config values also uses the names from the lookup_* tables
  // and include a full project name in the format <unit-short-name>-<stage-name>.
  projects_flattened = merge([
    for parent, items in var.projects: {
      for name, conf in items: "${parent}.${name}" => merge(
        conf, {
          parent  = local.stages_lookup[parent],
          name    = local.project_lookup[name],
          project = "${local.project_lookup[name]}-${local.stages_lookup[parent]}",
        }
      )
    }
  ]...)

  // all configured stage names as flat list
  stages = [for k in keys(var.projects) : local.stages_lookup[k]]
}

output stages {
  value = local.stages
}

output projects_flattened {
  value = local.projects_flattened
}
