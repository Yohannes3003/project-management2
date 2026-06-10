# Project Management API

Backend REST API untuk sistem Project Management menggunakan Golang, PostgreSQL, JWT Authentication, dan Gin Framework.

## Tech Stack

- Golang
- PostgreSQL
- Gin Framework
- JWT Authentication
- GORM
- Docker (Opsional)

## Fitur

- Authentication Login & Register
- JWT Authorization
- CRUD User
- CRUD Project
- CRUD Task
- Refresh Token
- Pagination
- Search & Filter

## Struktur Project

```
project-management-api/
├── cmd/
├── config/
├── controllers/
├── middleware/
├── models/
├── repositories/
├── routes/
├── services/
├── utils/
├── .env
├── go.mod
└── main.go
```

## Instalasi

### Clone Repository

```bash
git clone https://github.com/username/project-management-api.git

cd project-management-api
```

### Install Dependency

```bash
go mod tidy
```

### Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE project_management;
```

### Konfigurasi Environment

Buat file `.env`

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=passwordanda
DB_NAME=project_management

JWT_SECRET=supersecret
JWT_EXPIRED=1h
```

### Jalankan Project

```bash
go run main.go
```

atau

```bash
air
```

## API Endpoint

### Authentication

| Method | Endpoint | Description |
|----------|----------|-------------|
| POST | /api/register | Register User |
| POST | /api/login | Login User |

### Project

| Method | Endpoint |
|----------|----------|
| GET | /api/projects |
| POST | /api/projects |
| PUT | /api/projects/:id |
| DELETE | /api/projects/:id |

### Task

| Method | Endpoint |
|----------|----------|
| GET | /api/tasks |
| POST | /api/tasks |
| PUT | /api/tasks/:id |
| DELETE | /api/tasks/:id |

## Menjalankan Test

```bash
go test ./...
```

## Author

Yohannes Rahul
