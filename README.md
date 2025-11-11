## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose (for local development)

### Run
To run the project simply use docker-compose 
```bash
# Start all services
docker-compose up --build -d
```
or use make
```bash
# Start all services
make run
```
open http://localhost:8080/swagger

## Key Principles

- **No Direct Inter-Module Imports**: Modules communicate only through defined ports
- **Clean Architecture Layers**: Domain → Application → Infrastructure separation
- **Dependency Injection**: All dependencies injected via constructors
- **Pure Domain Layer**: Business logic with zero external dependencies

## 📁 Project Structure

```
.
├── api/                # API contracts (OpenAPI, protobuf)
├── cmd/
│   └── app/            # Main application entry point
├── resources/          # Configuration files, OpenApi, etc
├── internal/           # Private application code
│   └── [modules]/      # Business domain modules
├── platform/           # Shared infrastructure code
├── pkg/                # Public library code
├── scripts/            # Migration and scripts
├── docker-compose.yml  # Local development environment
├── Dockerfile          # Container image definition
└── Makefile            # Development commands
```

## 🧪 Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run all tests with coverage
make test-coverage
```