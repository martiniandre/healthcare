output "fhir_store" {
  value       = google_healthcare_fhir_store.main.name
  description = "FHIR store name"
}

output "dicom_store" {
  value       = google_healthcare_dicom_store.main.name
  description = "DICOM store name"
}

output "imaging_bucket" {
  value       = google_storage_bucket.dicom_imaging.name
  description = "GCS bucket for DICOM imaging"
}

output "postgres_connection_name" {
  value       = google_sql_database_instance.postgres.connection_name
  description = "Cloud SQL Postgres connection name"
}

output "redis_host" {
  value       = google_redis_instance.redis.host
  description = "Memorystore Redis host"
}

output "postgres_password" {
  value       = random_password.postgres_password.result
  sensitive   = true
  description = "Postgres application user password"
}
