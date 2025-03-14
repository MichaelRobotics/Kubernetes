# S3 Bucket for Velero Backups
resource "aws_s3_bucket" "velero" {
  bucket = "my-velero-bucket-v152"
  tags = {
    Name        = "Velero Backup Bucket"
    Environment = "Production"
  }
}

# S3 Bucket Versioning
resource "aws_s3_bucket_versioning" "velero_versioning" {
  bucket = aws_s3_bucket.velero.id

  versioning_configuration {
    status = "Enabled"
  }
}
