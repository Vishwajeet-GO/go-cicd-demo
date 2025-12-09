# Task Management API

![CI/CD Pipeline](https://github.com/Vishwajeet-GO/go-cicd-demo/workflows/Go%20CI/CD%20Pipeline/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED)
![License](https://img.shields.io/badge/license-MIT-green)

> RESTful API for task management built with Go, PostgreSQL, and Docker

## Features

- RESTful API with 5 CRUD endpoints
- PostgreSQL database integration
- Docker containerization with multi-stage builds
- Automated CI/CD pipeline with GitHub Actions
- Unit tests with coverage reporting
- Health monitoring endpoints

## Tech Stack

- **Language:** Go 1.25.5
- **Framework:** Gin
- **Database:** PostgreSQL 18
- **Containerization:** Docker, Docker Compose
- **CI/CD:** GitHub Actions
- **Testing:** Go testing package

## 🏗️ Architecture
```
go-cicd-demo/
├── cmd/api/              # Application entry point
├── internal/
│   ├── database/         # Database connection
│   ├── handlers/         # HTTP handlers
│   └── models/           # Data models
├── tests/                # Unit tests
├── Dockerfile            # Multi-stage Docker build
└── docker-compose.yml    # Local development setup
```

## Quick Start

### Prerequisites
- Go 1.25.5
- Docker and Docker compose
- PostgreSQL 18+ (for local development)

### Running with Docker Compose (Recommended)
```bash
# Clone repository
git clone https://github.com/Vishwajeet-GO/go-cicd-demo
cd go-cicd-demo

# Start all services
docker-compose up --build
```

### Running Locally
```bash
# Install dependencies
go mod download

# Set up environment variables 
cp .env.example .env


# Run application
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/tasks` | Create task |
| GET | `/api/tasks` | Get all tasks |
| GET | `/api/tasks/:id` | Get task by ID |
| PUT | `/api/tasks/:id` | Update task |
| DELETE | `/api/tasks/:id` | Delete task |
| GET | `/health` | Health check |


## Testing
```bash
# Run tests
go test ./tests-v

# Run tests with coverage
go test ./tests -cover

# Generate coverage report 
go test ./tests -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Docker

### Build Image
```bash
docker build -t go-task-api:v1 .
```

### Image Size Optimization

- **Before:** 350MB (standard Go image)
- **After:** 22MB (multi-stage Alpine build)
- **Reduction:** 93%

## Performance Metrics

- **Deployment Time Reduction:** 80% (15 min → 3 min)
- **Docker Image Size Reduction:** 93% (350MB → 22MB)
- **Test Coverage:** 35%+
- **Concurrent Connections:** 25 (connection pooling)

## Development

### Project Structure

- Clean architecture with separation of concerns
- Modular handlers, models, and database layers
- Environment-based configuration
- Graceful shutdown handling

### Code Quality

- Automated formatting checks (go fmt)
- Static analysis (go vet)
- Unit tests for critical paths
- CI/CD pipeline for every commit

## License

This project is open source and available under the MIT License.

## Author

**Vishwajeet Yadav**

- GitHub: @Vishwajeet-GO
- Email: yadavvishwajeet2004@gmail.com

## Acknowledgments

- Go community for excellent documentation
- Gin framework for elegant API development
- Docker for containerization simplicity
