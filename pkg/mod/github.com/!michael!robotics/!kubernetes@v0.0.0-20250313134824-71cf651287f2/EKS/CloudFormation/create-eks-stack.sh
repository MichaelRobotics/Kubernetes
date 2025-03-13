aws cloudformation create-stack \
  --region us-east-1 \
  --stack-name my-eks-cluster \
  --capabilities CAPABILITY_NAMED_IAM \
  --template-body file://eks-stack.yaml 