output "endpoint" {
  value = aws_eks_cluster.michaelrobotics.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.michaelrobotics.certificate_authority[0].data
}
output "cluster_id" {
  value = aws_eks_cluster.michaelrobotics.id
}
output "cluster_endpoint" {
  value = aws_eks_cluster.michaelrobotics.endpoint
}
output "cluster_name" {
  value = aws_eks_cluster.michaelrobotics.name
}