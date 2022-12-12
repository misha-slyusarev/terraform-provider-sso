terraform {
  required_providers {
    mysso = {
      version = "~> 0"
      source  = "mysso"
    }
  }
}

data "mysso_permission_set" "my_access" {
  name             = "MyAccess"
  description      = "Testing data source"
  session_duration = "D1H8"
}

resource "mysso_permission_set" "developer_access" {
  name             = "DeveloperAccess"
  description      = "General access for R&D"
  session_duration = data.mysso_permission_set.my_access.session_duration
  relay_state      = data.mysso_permission_set.my_access.rendered
}

resource "mysso_permission_set" "data_area_access" {
  name             = "DataAreaAccess"
  description      = "Access for the data area teams"
  session_duration = data.mysso_permission_set.my_access.session_duration
  relay_state      = data.mysso_permission_set.my_access.rendered
}

# --- Examples for the permission pool functionality

data "mysso_permission_pool" "sso_access" {
  relay_state = "ABC"
  tags        = { "test" = "tag" }

  permission_set {
    name               = "DataAreaAccess"
    description        = "Testing data source"
    session_duration   = "D1H8"
    relay_state        = "XYZ"
    policy_attachments = []
    based_on           = ["DeveloperAccess"]
  }

  permission_set {
    name             = "DeveloperAccess"
    description      = "Testing data source 2"
    session_duration = "D1H8"
    tags             = {
      "permission_set" = "permission tag"
    }
    policy_attachments = [
      "arn:aws:iam::aws:policy/DeveloperAccess",
    ]
  }
}

resource "mysso_permission_set" "pool_access" {
  for_each         = local.permissions
  name             = each.value.name
  description      = each.value.description
  session_duration = each.value.session_duration
  relay_state      = each.value.relay_state
  tags             = each.value.tags
}

locals {
  # Convert the permissions list to a map so it can be used with for_each
  permissions = {
    for p in data.mysso_permission_pool.sso_access.permissions : p.id => p
  }
  policy_attachments = {
    for p in data.mysso_permission_pool.sso_access.policy_attachments : p.id => p
  }
}
