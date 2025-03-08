# CI/CD Setup for OpenTelemetry Demo

This guide explains how to set up the CI/CD pipeline for building and pushing Docker images to your private Docker Hub repository.

## Repository Structure Note

This project is organized in a subdirectory structure:
- The OpenTelemetry Demo code is in the `opentelemetry-demo/` subdirectory
- The workflow watches for changes to files in both the root directory and the `opentelemetry-demo/` subdirectory

Make sure to modify files in the correct location to trigger the workflow.

## Required GitHub Secrets

Before you can use the workflow, you need to set up the following secrets in your GitHub repository:

1. `DOCKER_USERNAME`: Your Docker Hub username
2. `DOCKER_PASSWORD`: Your Docker Hub password or access token (recommended)

## Setting Up GitHub Secrets

1. Go to your GitHub repository
2. Click on "Settings" tab
3. In the left sidebar, click on "Secrets and variables" → "Actions"
4. Click on "New repository secret"
5. Add each of the required secrets:
   - Name: `DOCKER_USERNAME`
   - Value: Your Docker Hub username
   - Click "Add secret"
   - Name: `DOCKER_PASSWORD`
   - Value: Your Docker Hub password or access token
   - Click "Add secret"

## About the CI/CD Workflow

The workflow (`custom-ci-cd.yml`) does the following:

1. Triggers on:
   - Push to the `main` branch (when changes are made to `src/`, `test/`, or `.env.override`)
   - Pull requests to the `main` branch
   - Manual trigger (workflow_dispatch) with a version parameter

2. Jobs:
   - `protobufcheck`: Ensures Protocol Buffers are correctly generated
   - `build_and_push_images`: Builds and pushes Docker images for all microservices to `robclusterdev/clusterimages`
   - `deploy_demo`: Updates the `.env.override` file for deployment

3. The images are tagged as:
   - `robclusterdev/clusterimages:<service-name>`
   - `robclusterdev/clusterimages:<version>-<service-name>` (version is either "latest" or specified during manual trigger)

## Running the Workflow Manually

To run the workflow manually with a custom version:

1. Go to the "Actions" tab in your repository
2. Click on "Custom CI/CD Pipeline" workflow
3. Click "Run workflow"
4. Enter a version tag (e.g., "v1.0.0")
5. Click "Run workflow"

## Local Testing

To test locally:

1. Ensure Docker is installed and you're logged in to Docker Hub
2. Pull your repository
3. Run `make start` to start the application with your custom images 