# OpenTelemetry Demo

## 🛠️ Makefile Guide

The Makefile serves as the main interface for working with this project.
Below is a detailed breakdown of available commands:

### Setup and Configuration

The Makefile automates Docker Compose commands with proper environment configuration:

```make
DOCKER_COMPOSE_CMD ?= docker compose
DOCKER_COMPOSE_ENV=--env-file .env --env-file .env.override
```

### Core Commands

#### Application Lifecycle

| Command | Description |
|---------|-------------|
| `make build` | Builds all containers |
| `make start` | Starts the demo application in detached mode |
| `make start-minimal` | Starts a minimal version of the demo |
| `make stop` | Stops all containers and removes volumes |
| `make restart service=X` | Restarts a specific service |
| `make redeploy service=X` | Rebuilds and restarts a specific service |

#### Testing

| Command | Description |
|---------|-------------|
| `make run-tests` | Runs all tests (frontend and trace-based) |
| `make run-tracetesting` | Runs trace-based tests (optionally for specific services) |

#### Code Generation

| Command | Description |
|---------|-------------|
| `make generate-protobuf` | Generates protocol buffer code using local tools |
| `make docker-generate-protobuf` | Generates protocol buffer code using Docker (recommended) |
| `make clean` | Removes generated protobuf files |
| `make check-clean-work-tree` | Verifies no uncommitted changes exist |

#### Docker Image Management

| Command | Description |
|---------|-------------|
| `make build-and-push` | Builds and pushes images to a registry |
| `make build-multiplatform` | Builds images for multiple architectures (AMD64, ARM64) |
| `make clean-images` | Removes OpenTelemetry Demo images |

#### Code Quality

| Command | Description |
|---------|-------------|
| `make check` | Runs all code quality checks |
| `make misspell` | Checks for spelling errors in documentation |
| `make markdownlint` | Validates Markdown formatting |
| `make yamllint` | Validates YAML files |
| `make checklicense` | Ensures license headers are present |

### Example Usage

```bash
# Build and start the demo
make build
make start
# Run tests
make run-tests
# Rebuild and restart a specific service
make redeploy service=frontend
# Stop the demo
make stop
```
