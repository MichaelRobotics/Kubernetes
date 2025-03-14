# Output Access Key and Secret
output "velero_access_key_id" {
  value       = aws_iam_access_key.velero_access_key.id
  description = "Velero Access Key ID"
}

output "velero_secret_access_key" {
  value       = aws_iam_access_key.velero_access_key.secret
  description = "Velero Secret Access Key"
  sensitive   = true
}