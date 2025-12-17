github_token              = ""
github_owner              = "GDJ4"
repository_name           = "iac-lab-terraform"
repository_description    = "Практическая работа 7.1: Terraform + GitHub"
repository_visibility     = "public"
default_branch            = "main"
license_template          = "mit"
repository_topics         = ["terraform", "iac", "github"]
enable_projects           = false
allow_merge_commit        = true
allow_squash_merge        = true
allow_rebase_merge        = true
allow_auto_merge          = false
delete_branch_on_merge    = true
enable_vulnerability_alerts = true
require_signed_commits    = true
require_code_owner_reviews = false
required_review_count     = 1
required_status_checks    = ["ci/build"]
enforce_admins            = false

collaborators = [
  # { username = "teammate", permission = "push" },
]

webhooks = [
  # {
  #   url          = "https://example.com/webhook"
  #   events       = ["push", "pull_request"]
  #   content_type = "json"
  #   active       = true
  # }
]

codeowners = [
  # { path = "*", owners = ["@your-github-username"] },
]
