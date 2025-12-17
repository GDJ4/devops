terraform {
  required_version = ">= 1.6.0"

  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
}

provider "github" {
  token = var.github_token
  owner = var.github_owner
}

resource "github_repository" "iac_repo" {
  name                  = var.repository_name
  description           = var.repository_description
  visibility            = var.repository_visibility
  auto_init             = true
  has_issues            = true
  has_wiki              = false
  has_projects          = var.enable_projects
  allow_merge_commit    = var.allow_merge_commit
  allow_squash_merge    = var.allow_squash_merge
  allow_rebase_merge    = var.allow_rebase_merge
  allow_auto_merge      = var.allow_auto_merge
  delete_branch_on_merge = var.delete_branch_on_merge
  license_template      = var.license_template
  vulnerability_alerts  = var.enable_vulnerability_alerts
  topics                = var.repository_topics

  lifecycle {
    prevent_destroy = true
  }
}

resource "github_branch_default" "default" {
  repository = github_repository.iac_repo.name
  branch     = var.default_branch
}

resource "github_branch_protection" "default" {
  repository_id = github_repository.iac_repo.node_id
  pattern       = var.default_branch

  depends_on = [
    github_branch_default.default,
    github_repository_file.codeowners
  ]

  enforce_admins                = var.enforce_admins
  allows_deletions              = false
  allows_force_pushes           = false
  require_conversation_resolution = true
  required_linear_history       = true
  require_signed_commits        = var.require_signed_commits

  dynamic "required_status_checks" {
    for_each = length(var.required_status_checks) == 0 ? [] : [var.required_status_checks]

    content {
      strict   = true
      contexts = required_status_checks.value
    }
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = true
    require_code_owner_reviews      = var.require_code_owner_reviews
    required_approving_review_count = var.required_review_count
  }
}

resource "github_repository_collaborator" "repo_collaborators" {
  for_each = { for collaborator in var.collaborators : collaborator.username => collaborator }

  repository = github_repository.iac_repo.name
  username   = each.value.username
  permission = each.value.permission
}

resource "github_repository_webhook" "repo_webhooks" {
  for_each = { for webhook in var.webhooks : webhook.url => webhook }

  repository = github_repository.iac_repo.name
  active     = each.value.active
  events     = each.value.events

  configuration {
    url          = each.value.url
    content_type = each.value.content_type
    insecure_ssl = false
    secret       = try(each.value.secret, null)
  }
}

resource "github_repository_file" "codeowners" {
  count = length(var.codeowners) == 0 ? 0 : 1

  repository     = github_repository.iac_repo.name
  branch         = var.default_branch
  file           = ".github/CODEOWNERS"
  commit_message = "Add CODEOWNERS managed by Terraform"
  overwrite_on_create = true

  content = join("\n", concat(
    ["# Managed by Terraform. Edit in IaC."],
    [for rule in var.codeowners : "${rule.path} ${join(" ", rule.owners)}"]
  ))

  depends_on = [github_branch_default.default]
}
