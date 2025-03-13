resource "aws_iam_role" "demo2" {
  name = "eks-cluster-demo2"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "demo2_amazon_eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.demo2.name
}

resource "aws_eks_cluster" "demo2" {
  name     = "demo2"
  role_arn = aws_iam_role.demo2.arn

  vpc_config {
    subnet_ids = [
      aws_subnet.private_us_east_2a.id,
      aws_subnet.private_us_east_2b.id,
      aws_subnet.public_us_east_2a.id,
      aws_subnet.public_us_east_2b.id
    ]
  }

  depends_on = [aws_iam_role_policy_attachment.demo2_amazon_eks_cluster_policy]
}