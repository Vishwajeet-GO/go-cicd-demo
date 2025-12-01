# Task Management API

RESTful API for task management built with Go and PostgreSQL.

## Features

- RESTful API with CRUD operations
- PostgreSQL database integration
- Health check endpoints
- Graceful shutdown handling

## Tech Stack

- Go 1.21+
- Gin Framework
- PostgreSQL 15+

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+

### Installation
```bash
git clone https://github.com/Vishwajeet-GO/go-cicd-demo.git
cd go-cicd-demo
go mod download
cp .env.example .env
# Edit .env with your database credentials
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/tasks | Create task |
| GET | /api/tasks | Get all tasks |
| GET | /api/tasks/:id | Get task by ID |
| PUT | /api/tasks/:id | Update task |
| DELETE | /api/tasks/:id | Delete task |
| GET | /health | Health check |

## Author

**Vishwajeet Yadav**
- GitHub: @Vishwajeet-GO
- Email: yadavvishwajeet@gmail.com
