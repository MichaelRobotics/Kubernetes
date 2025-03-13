# Kubernetes Learning Repository

This repository contains comprehensive Kubernetes resources, configurations, and examples for both learning and real-world implementation.

## 🌟 Repository Highlights

- **Infrastructure**: Configurations for AKS, EKS, Kind, Minikube
- **Security**: mTLS, cert-manager, Kyverno, Gatekeeper examples
- **Networking**: Ingress controllers, network policies, Calico
- **Storage**: CSI drivers, persistent volumes, database deployments
- **Advanced Features**: Autoscaling (HPA, VPA), StatefulSets, more
- **DevOps**: CI/CD pipelines, infrastructure as code with Terraform
- **Service Mesh**: Linkerd configuration examples

## 📚 Learning Paths

1. **Beginners**: Start with basic workloads (Deployments, Services)
2. **Cloud Engineers**: Explore cloud-specific implementations
3. **Security Professionals**: Study security-focused sections
4. **DevOps Engineers**: Examine CI/CD and automation examples
5. **SREs**: Review monitoring, scaling, and resilience patterns

## 🔧 Tools & Technologies

- **Container Orchestration**: Kubernetes
- **Cloud Providers**: AWS (EKS), Azure (AKS)
- **Infrastructure as Code**: Terraform
- **CI/CD**: Jenkins
- **Observability**: EFK stack
- **Packaging**: Helm, Kustomize
- **Policy Enforcement**: Kyverno, Gatekeeper

## 🚀 For Recruiters

This repository demonstrates expertise in cloud-native architectures, security best practices, automation, and solving real-world Kubernetes challenges with production-ready implementations.

## 🔄 GitHub Workflows

- **custom-ci-cd.yml**: Main CI/CD pipeline for building and deploying
- **component-build-images.yml**: Reusable workflow for container image building
- **build-images.yml**: Quick verification of container images
- **checks.yml**: Code quality verification
- **run-integration-tests.yml**: Executes end-to-end tests
- **dependabot-auto-update-protobuf-diff.yml**: Maintains protobuf files
- **gradle-wrapper-validation.yml**: Security check for Gradle wrappers

## Security Notice

**IMPORTANT**: This repository is primarily for educational purposes. Always review and secure credentials before using any configurations in production environments.

