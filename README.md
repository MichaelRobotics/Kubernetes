# Kubernetes Learning Repository

This repository is a comprehensive collection of Kubernetes resources, configurations, and examples for both learning and real-world implementation. It covers a wide range of Kubernetes concepts and integrations with various cloud platforms and tools.

## 🌟 Repository Highlights

- **Diverse Infrastructure Setups**: Configurations for AKS, EKS, Kind, Minikube, and more
- **Security Implementations**: mTLS, cert-manager, Kyverno, and Gatekeeper examples
- **Networking Solutions**: Ingress controllers, network policies, and Calico configurations
- **Storage Options**: CSI drivers, persistent volumes, and database deployments
- **Advanced Kubernetes Features**: Autoscaling (HPA, VPA, Cluster Autoscaler), StatefulSets, and more
- **DevOps Practices**: CI/CD pipelines, infrastructure as code with Terraform
- **Service Mesh**: Linkerd configuration examples

## 📚 Learning Paths

This repository can serve as a learning path for:

1. **Kubernetes Beginners**: Start with basic workloads (Deployments, Services)
2. **Cloud Engineers**: Explore cloud-specific implementations (AKS, EKS)
3. **Security Professionals**: Study the security-focused sections (mTLS, Kyverno)
4. **DevOps Engineers**: Examine the CI/CD and automation examples
5. **SREs**: Review monitoring, scaling, and resilience patterns

## 🔧 Tools & Technologies

- **Container Orchestration**: Kubernetes
- **Cloud Providers**: AWS (EKS), Azure (AKS)
- **Infrastructure as Code**: Terraform
- **CI/CD**: Jenkins
- **Observability**: EFK stack
- **Packaging**: Helm, Kustomize
- **Policy Enforcement**: Kyverno, Gatekeeper
- **Controllers**: Metacontroller, Custom controllers

## 🚀 For Recruiters

This repository demonstrates expertise in:

- **Modern Infrastructure Design**: Cloud-native architectures and implementation
- **Security-First Approach**: Implementation of security best practices
- **Automation Skills**: Infrastructure as code and automation workflows
- **Problem-Solving**: Real-world solutions to common Kubernetes challenges
- **Technical Breadth**: Knowledge across multiple technologies and platforms
- **Production-Ready Implementations**: Enterprise patterns applicable in the real world

Feel free to explore the different directories to see practical examples of Kubernetes deployments, configurations, and integrations.

## 🔄 GitHub Workflows

The repository contains several optimized GitHub workflows for CI/CD automation. Each workflow has been configured to trigger at the appropriate time without overlapping:

### Core Workflows

1. **custom-ci-cd.yml** - PRIMARY WORKFLOW
   - Purpose: Main CI/CD pipeline for building and deploying the opentelemetry-demo
   - Triggers: Only on push to main branch or manual workflow dispatch
   - Functions: Builds and pushes Docker images, performs protobuf checks, deploys the demo
   
2. **component-build-images.yml**
   - Purpose: Reusable workflow that handles container image building
   - Usage: Not triggered directly - called by other workflows
   - Functions: Provides Docker build functionality for other workflows
   
3. **build-images.yml**
   - Purpose: Quick verification of container images
   - Triggers: Pull requests affecting src/**, test/**, or Dockerfile files
   - Functions: Validates builds without pushing images
   
4. **checks.yml**
   - Purpose: Code quality checks
   - Triggers: Pull requests and pushes to non-main branches
   - Excludes: Markdown files and documentation changes
   - Functions: Runs markdownlint and validates images without pushing

### Supporting Workflows

5. **run-integration-tests.yml**
   - Purpose: Runs integration tests
   - Triggers: On pull request reviews and pull requests
   - Functions: Executes end-to-end tests with tracetesting
   
6. **dependabot-auto-update-protobuf-diff.yml**
   - Purpose: Maintains protobuf files when dependencies change
   - Triggers: PRs affecting dependency files (package.json, go.mod, etc.)
   - Functions: Runs only for dependabot PRs to auto-update protobuf files
   
7. **gradle-wrapper-validation.yml**
   - Purpose: Security check for Gradle wrappers
   - Triggers: Only when Gradle wrapper files change
   - Functions: Validates wrapper checksums for security

8. **label-pr.yml**
   - Purpose: Automatically labels pull requests
   - Triggers: When PRs are opened, synchronized, or reopened
   - Functions: Adds labels based on changed files for easier PR management

## Security Notice

**IMPORTANT**: This repository is primarily for educational purposes. Always review and secure credentials before using any configurations in production environments.

