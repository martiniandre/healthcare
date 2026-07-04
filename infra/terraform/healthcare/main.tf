resource "google_project_service" "healthcare_api" {
  service = "healthcare.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudsql_admin" {
  service = "sqladmin.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "redis_api" {
  service = "redis.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "storage_api" {
  service = "storage.googleapis.com"
  disable_on_destroy = false
}

resource "google_healthcare_dataset" "main" {
  name     = "${var.environment}-healthcare-dataset"
  location = var.region
  project  = var.project_id
  depends_on = [google_project_service.healthcare_api]
}

resource "google_healthcare_fhir_store" "main" {
  name    = "${var.environment}-fhir-store"
  dataset = google_healthcare_dataset.main.id
  version = "R4"
  depends_on = [google_project_service.healthcare_api]
}

resource "google_healthcare_dicom_store" "main" {
  name    = "${var.environment}-dicom-store"
  dataset = google_healthcare_dataset.main.id
  depends_on = [google_project_service.healthcare_api]
}

resource "google_storage_bucket" "dicom_imaging" {
  name     = "${var.project_id}-${var.environment}-imaging"
  location = var.region
  uniform_bucket_level_access = true
  versioning {
    enabled = true
  }
  lifecycle_rule {
    condition {
      age = 365
    }
    action {
      type = "Delete"
    }
  }
  depends_on = [google_project_service.storage_api]
}

resource "google_sql_database_instance" "postgres" {
  name             = "${var.environment}-healthcare-pg"
  database_version = "POSTGRES_16"
  region           = var.region
  settings {
    tier              = "db-custom-2-7680"
    availability_type = "REGIONAL"
    disk_size         = 100
    disk_type         = "PD_SSD"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "03:00"
      transaction_log_retention_days = 7
    }
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_service_networking_connection.private_vpc_connection.network
    }
  }
  depends_on = [google_project_service.cloudsql_admin]
}

resource "google_sql_database" "healthcare" {
  name     = "healthcare"
  instance = google_sql_database_instance.postgres.name
}

resource "google_sql_user" "postgres_app" {
  name     = "app"
  instance = google_sql_database_instance.postgres.name
  password = random_password.postgres_password.result
}

resource "random_password" "postgres_password" {
  length  = 24
  special = false
}

resource "google_redis_instance" "redis" {
  name           = "${var.environment}-healthcare-redis"
  tier           = "STANDARD_HA"
  memory_size_gb = 2
  region         = var.region
  redis_version  = "REDIS_7"
  connect_mode   = "PRIVATE_SERVICE_ACCESS"
  depends_on     = [google_project_service.redis_api]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          = "${var.environment}-private-ip-alloc"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

data "google_compute_network" "default" {
  name = "default"
}

output "fhir_store" {
  value = google_healthcare_fhir_store.main.name
}

output "dicom_store" {
  value = google_healthcare_dicom_store.main.name
}

output "imaging_bucket" {
  value = google_storage_bucket.dicom_imaging.name
}

output "postgres_connection_name" {
  value = google_sql_database_instance.postgres.connection_name
}

output "redis_host" {
  value = google_redis_instance.redis.host
}

output "postgres_password" {
  value     = random_password.postgres_password.result
  sensitive = true
}
