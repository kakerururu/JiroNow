resource "google_artifact_registry_repository" "backend_repo" {
  provider = google
  location = "asia-northeast1"
  repository_id = "backend"
  format = "DOCKER"
  description = "Backend Docker repository"
}

resource "google_iam_workload_identity_pool" "github_pool" {
  workload_identity_pool_id = "github"
  display_name              = "GitHub"
  description               = "GitHub Actions 用 Workload Identity Pool"
  disabled                  = false
}

resource "google_iam_workload_identity_pool_provider" "github_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub Actions Provider"
  description                        = "OIDC Provider for GitHub"
  attribute_mapping = {
    "google.subject" = "assertion.sub"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account" "github_sa" {
  account_id   = "github-actions-deploy er"
  display_name = "GitHub Actions SA"
}

resource "google_service_account_iam_binding" "github_wif_binding" {
  service_account_id = google_service_account.github_sa.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/YOUR_GITHUB_REPO"
  ]
}
