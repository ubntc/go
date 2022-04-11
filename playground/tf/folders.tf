module "stage_folders" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 3.0"

  parent = var.organization_id
  names  = local.stages

  set_roles = false
}
