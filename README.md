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

The repository contains several GitHub workflows for CI/CD automation. Here's an overview of the key workflows:

### Core Workflows

1. **custom-ci-cd.yml** - PRIMARY WORKFLOW
   - Purpose: Main CI/CD pipeline for building and deploying the opentelemetry-demo
   - Triggers: Push to main branch affecting opentelemetry-demo directory
   
2. **component-build-images.yml**
   - Purpose: Builds container images for individual components
   - Used by: custom-ci-cd and other workflows

3. **build-images.yml**
   - Purpose: Tests image generation for changes to src/test files
   - Triggers: Push events affecting source or test files

4. **checks.yml**
   - Purpose: Runs linting and other code quality checks
   - Triggers: Push and pull requests to main branch

### Additional Workflows (Consider Removing)

1. **stale.yml**
   - Purpose: Automatically closes stale PRs
   - Recommended: Only keep if you have many contributors and PRs

2. **nightly-release.yml**
   - Purpose: Creates nightly builds of all components
   - Recommended: Only keep if you need nightly builds for testing

3. **release.yml**
   - Purpose: Publishes container images for GitHub releases
   - Recommended: Only keep if you're publishing official releases

## Security Notice

**IMPORTANT**: This repository is primarily for educational purposes. Always review and secure credentials before using any configurations in production environments.

