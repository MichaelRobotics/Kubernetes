resource "aws_iam_role" "nodes2" {
  name = "eks-node-group-nodes2"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "nodes_amazon_eks_worker_node_policy2" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.nodes2.name
}

resource "aws_iam_role_policy_attachment" "nodes_amazon_eks_cni_policy2" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.nodes2.name
}

resource "aws_iam_role_policy_attachment" "nodes_amazon_ec2_container_registry_read_only2" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.nodes2.name
}

resource "aws_eks_node_group" "private_nodes2" {
  cluster_name    = aws_eks_cluster.demo2.name
  node_group_name = "private-nodes2"
  node_role_arn   = aws_iam_role.nodes2.arn

  subnet_ids = [
    aws_subnet.private_us_east_2a.id,
    aws_subnet.private_us_east_2b.id
  ]

  capacity_type  = "ON_DEMAND"
  instance_types = ["t3.medium"]

  scaling_config {
    desired_size = 1
    max_size     = 5
    min_size     = 0
  }

  update_config {
    max_unavailable = 1
  }

  labels = {
    role = "general"
  }

  depends_on = [
    aws_iam_role_policy_attachment.nodes_amazon_eks_worker_node_policy2,
    aws_iam_role_policy_attachment.nodes_amazon_eks_cni_policy2,
    aws_iam_role_policy_attachment.nodes_amazon_ec2_container_registry_read_only2,
  ]
}