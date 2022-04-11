
terraform {
  required_version = ">= 1.0.0"
  # backend "gcs" {
  #   bucket = "my-platform-tfstate-bucket"
  #   prefix = "platform"
  # }
}

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

output projects_flattened {
  value = local.projects_flattened
}
