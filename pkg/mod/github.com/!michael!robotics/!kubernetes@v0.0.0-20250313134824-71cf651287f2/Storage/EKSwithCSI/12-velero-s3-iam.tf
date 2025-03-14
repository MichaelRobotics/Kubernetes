resource "aws_iam_policy" "velero_policy" {
  name        = "VeleroPolicy"
  description = "IAM policy for Velero to restore EKS cluster and manage AWS resources"
  policy      = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = [
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = [
          aws_s3_bucket.velero.arn,
          "${aws_s3_bucket.velero.arn}/*"
        ]
      },
      {
        Effect   = "Allow"
        Action   = [
          "ec2:CreateSnapshot",
          "ec2:DeleteSnapshot",
          "ec2:DescribeSnapshots",
          "ec2:CreateTags",
          "ec2:DeleteTags",
          "ec2:AttachVolume",
          "ec2:DetachVolume",
          "ec2:DescribeVolumes",
          "ec2:ModifyVolume",
          "ec2:DeleteVolume",
          "ec2:CreateVolume"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "iam:PassRole"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:GenerateDataKey"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "eks:DescribeCluster",
          "eks:ListClusters",
          "eks:UpdateClusterConfig"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "ec2:DescribeVpcs",
          "ec2:DescribeSubnets",
          "ec2:DescribeSecurityGroups"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "elasticloadbalancing:DescribeLoadBalancers",
          "elasticloadbalancing:CreateLoadBalancer",
          "elasticloadbalancing:DeleteLoadBalancer",
          "elasticloadbalancing:ModifyLoadBalancerAttributes"
        ]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = [
          "elasticfilesystem:DescribeFileSystems",
          "elasticfilesystem:CreateFileSystem",
          "elasticfilesystem:DeleteFileSystem",
          "elasticfilesystem:CreateMountTarget"
        ]
        Resource = "*"
      }
    ]
  })
}


# IAM User for Velero
resource "aws_iam_user" "velero_user" {
  name = "velero-user"
  tags = {
    Name        = "Velero User"
    Environment = "Production"
  }
}

# Attach Policy to the IAM User
resource "aws_iam_user_policy_attachment" "velero_attachment" {
  user       = aws_iam_user.velero_user.name
  policy_arn = aws_iam_policy.velero_policy.arn
}

# Access Key for the IAM User
resource "aws_iam_access_key" "velero_access_key" {
  user = aws_iam_user.velero_user.name
}

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
output "velero_policy_arn" {
  value = aws_iam_policy.velero_policy.arn
}