output "fhir_store" {
  value       = module.healthcare.fhir_store
  description = "FHIR store name"
}

output "dicom_store" {
  value       = module.healthcare.dicom_store
  description = "DICOM store name"
}

output "imaging_bucket" {
  value       = module.healthcare.imaging_bucket
  description = "GCS bucket for DICOM imaging"
}

output "postgres_connection_name" {
  value       = module.healthcare.postgres_connection_name
  description = "Cloud SQL Postgres connection name"
}

output "redis_host" {
  value       = module.healthcare.redis_host
  description = "Memorystore Redis host"
}

output "postgres_password" {
  value       = module.healthcare.postgres_password
  sensitive   = true
  description = "Postgres application user password"
}
