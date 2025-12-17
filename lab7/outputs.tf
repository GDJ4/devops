output "repository_urls" {
  description = "Основные ссылки на созданный репозиторий."
  value = {
    html = github_repository.iac_repo.html_url
    ssh  = github_repository.iac_repo.ssh_clone_url
    http = github_repository.iac_repo.http_clone_url
  }
}

output "default_branch" {
  description = "Ветка по умолчанию."
  value       = github_branch_default.default.branch
}

output "collaborators" {
  description = "Коллабораторы и их права."
  value       = { for name, collab in github_repository_collaborator.repo_collaborators : name => collab.permission }
}

output "webhooks" {
  description = "Активные вебхуки и события."
  value       = { for url, hook in github_repository_webhook.repo_webhooks : url => hook.events }
}
