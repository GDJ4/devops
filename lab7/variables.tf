variable "github_token" {
  description = "Персональный токен GitHub с правами на управление репозиториями."
  type        = string
  sensitive   = true
}

variable "github_owner" {
  description = "Имя пользователя или организации GitHub, на которую создается репозиторий."
  type        = string
}

variable "repository_name" {
  description = "Имя создаваемого репозитория."
  type        = string
}

variable "repository_description" {
  description = "Описание репозитория."
  type        = string
  default     = "Terraform-managed GitHub repository for IaC lab."
}

variable "repository_visibility" {
  description = "Приватность репозитория: public или private."
  type        = string
  default     = "public"

  validation {
    condition     = contains(["public", "private"], var.repository_visibility)
    error_message = "Допустимые значения для repository_visibility: public или private."
  }
}

variable "license_template" {
  description = "Шаблон лицензии GitHub (например, mit, apache-2.0)."
  type        = string
  default     = "mit"
}

variable "default_branch" {
  description = "Ветка по умолчанию."
  type        = string
  default     = "main"
}

variable "repository_topics" {
  description = "Список тем (topics) репозитория."
  type        = list(string)
  default     = []
}

variable "enable_projects" {
  description = "Включить Projects в репозитории."
  type        = bool
  default     = false
}

variable "allow_merge_commit" {
  description = "Разрешить merge commit."
  type        = bool
  default     = true
}

variable "allow_squash_merge" {
  description = "Разрешить squash merge."
  type        = bool
  default     = true
}

variable "allow_rebase_merge" {
  description = "Разрешить rebase merge."
  type        = bool
  default     = true
}

variable "allow_auto_merge" {
  description = "Разрешить автоматическое слияние (auto-merge)."
  type        = bool
  default     = false
}

variable "delete_branch_on_merge" {
  description = "Удалять ветку после слияния Pull Request."
  type        = bool
  default     = true
}

variable "enable_vulnerability_alerts" {
  description = "Включить Dependabot alerts."
  type        = bool
  default     = true
}

variable "require_signed_commits" {
  description = "Требовать подписанные коммиты на защищенной ветке."
  type        = bool
  default     = true
}

variable "require_code_owner_reviews" {
  description = "Требовать ревью code owners."
  type        = bool
  default     = false
}

variable "required_review_count" {
  description = "Минимум approving reviews перед слиянием."
  type        = number
  default     = 1

  validation {
    condition     = var.required_review_count >= 1 && var.required_review_count <= 6
    error_message = "Количество обязательных ревью должно быть от 1 до 6."
  }
}

variable "required_status_checks" {
  description = "Список обязательных статус-чеков (например, ci/build)."
  type        = list(string)
  default     = []
}

variable "enforce_admins" {
  description = "Применять правила защиты веток к администраторам."
  type        = bool
  default     = false
}

variable "collaborators" {
  description = "Список коллабораторов и их прав доступа."
  type = list(object({
    username   = string
    permission = optional(string, "push")
  }))
  default = []
}

variable "webhooks" {
  description = "Настройки вебхуков репозитория."
  type = list(object({
    url          = string
    events       = list(string)
    content_type = optional(string, "json")
    secret       = optional(string)
    active       = optional(bool, true)
  }))
  default = []
}

variable "codeowners" {
  description = "Правила CODEOWNERS для защиты файлов и директорий."
  type = list(object({
    path   = string
    owners = list(string)
  }))
  default = []
}
