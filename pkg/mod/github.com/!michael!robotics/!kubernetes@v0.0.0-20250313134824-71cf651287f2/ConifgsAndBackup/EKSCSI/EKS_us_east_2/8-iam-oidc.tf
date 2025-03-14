data "tls_certificate" "eks2" {
  url = aws_eks_cluster.demo2.identity[0].oidc[0].issuer
}

resource "aws_iam_openid_connect_provider" "eks2" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.eks2.certificates[0].sha1_fingerprint]
  url             = aws_eks_cluster.demo2.identity[0].oidc[0].issuer
}