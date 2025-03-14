## Use the existing IAM policy (VeleroPolicy)
#data "aws_iam_policy" "velero_policy" {
#  arn = "arn:aws:iam::533267337200:policy/VeleroPolicy"  # Replace with the correct ARN of the existing policy
#}


# Use the existing IAM user (velero-user)
#data "aws_iam_user" "velero_user" {
#  user_name = "velero-user"  # Replace with the correct name of the existing user
#}

### Attach Policy to the existing IAM User
#resource "aws_iam_user_policy_attachment" "velero_attachment" {
#  user       = data.aws_iam_user.velero_user.user_name  # Correct reference to the imported IAM user
#  policy_arn = data.aws_iam_policy.velero_policy.arn  # Correct reference to the imported IAM policy
#}


# Create IAM Access Key for the imported user
#resource "aws_iam_access_key" "velero_access_key" {
#  user = data.aws_iam_user.velero_user.user_name  # Use the data source reference
#}
#
## Output Access Key and Secret
#output "velero_access_key_id" {
#  value       = aws_iam_access_key.velero_access_key.id
#  description = "Velero Access Key ID"
#}
#
#output "velero_secret_access_key" {
#  value       = aws_iam_access_key.velero_access_key.secret
#  description = "Velero Secret Access Key"
#  sensitive   = true
#}