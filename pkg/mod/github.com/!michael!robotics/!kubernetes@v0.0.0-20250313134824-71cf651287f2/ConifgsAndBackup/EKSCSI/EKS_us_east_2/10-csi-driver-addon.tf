resource "aws_eks_addon" "csi_driver" {
  cluster_name = aws_eks_cluster.demo2.name
  addon_name   = "aws-ebs-csi-driver"

  service_account_role_arn = aws_iam_role.eks_ebs_csi_driver2.arn  # Reference to the existing IAM role

  # Additional configurations for the addon
}