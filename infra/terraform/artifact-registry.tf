resource "google_artifact_registry_repository" "backend" {
  provider     = google
  location     = "asia-northeast1"
  repository_id = "backend"
  description  = "Docker images for backend"
  format       = "DOCKER"
}
