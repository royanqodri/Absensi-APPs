# Absensi API 🏢
### Employee Attendance Backend Service (Go + Echo + MySQL)

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![Echo](https://img.shields.io/badge/Echo-v4-00ADD8?style=for-the-badge)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-ORM-blue?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

A backend REST API for employee attendance management, built with **Go**, **Echo**, and **GORM/MySQL**, following a clean **Controller → Service → Repository** architecture. Built as a portfolio project to demonstrate backend fundamentals: structured layering, bulk operations, auto-migration, and containerized deployment.

---

## 📌 Table of Contents
- [About](#-about)
- [Key Features](#-key-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Deployment](#-deployment)
- [Roadmap](#-roadmap)
- [Author](#-author)

---

## 📌 About

This service manages two core resources: **employees (users)** and **attendance records**, including regular check-in/check-out times and overtime tracking. It's structured to be easy to extend — the repository/service split makes it straightforward to add authentication, caching, or new endpoints without touching business logic.

---

## ✨ Key Features

### Attendance & User Management
- ✅ CRUD for employees (`/api/v1/user`) with filtering & pagination
- ✅ CRUD for attendance records (`/api/v1/absensi`) — check-in time, check-out time, overtime in/out
- ✅ **Bulk create/update** — send an array of records in a single request
- ✅ Duplicate-entry detection with structured error responses
- ✅ Password hashing with MD5 + bcrypt on create/update

### Data Layer
- ✅ MySQL via GORM with **auto-migration on startup**
- ✅ Connection pooling (100 max open, 10 idle, 1-hour max lifetime)
- ✅ Soft delete support (`deleted_at`) on all entities

### Operations
- ✅ `/health` endpoint with live DB connectivity check
- ✅ Graceful shutdown on SIGTERM/SIGINT
- ✅ Structured logging with Logrus
- ✅ Dockerized with an automated build → push → SSH-deploy pipeline (GitHub Actions)

### Prepared, not yet wired in
- 🔧 JWT signing/verification helpers (`middlewares/jwt-middleware.go`) exist but aren't attached to any route yet — no login endpoint issues tokens today
- 🔧 Redis helper functions exist (`util/redis.go`) but the client isn't initialized — caching isn't active yet

---

## 🛠️ Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.26+ |
| Web Framework | Echo v4 |
| ORM | GORM |
| Database | MySQL 8.0 |
| Password Hashing | bcrypt + MD5 |
| Configuration | Viper |
| Logging | Logrus |
| Containerization | Docker |
| CI/CD | GitHub Actions (build, push to Docker Hub, deploy via SSH) |

---

## 🏗️ Architecture

```
Client → Controller (request validation, HTTP response)
       → Service (business logic, transactions)
       → Repository (GORM queries)
       → MySQL
```

Each layer only talks to the one directly below it, which keeps business logic independent of both the HTTP framework and the database driver — useful if the API layer or the database ever needs to change.

---

## 📁 Project Structure

```
Absensi-APPs/
├── cmd/main.go              # Entry point, router setup, graceful shutdown
├── config/                  # Env-based configuration (Viper)
├── controller/               # HTTP handlers (Echo)
├── service/                  # Business logic
├── repository/                # GORM data access
├── model/
│   ├── entity/                # DB models (TUser, TAbsensi)
│   ├── request/                # Request DTOs
│   └── response/               # Response DTOs
├── middlewares/               # JWT helper (not yet wired to routes)
├── database/                  # MySQL connection + auto-migration
├── util/                      # Hashing, Redis helpers, response formatting, etc.
├── .github/workflows/deploy.yml  # Build → push → deploy pipeline
├── Dockerfile
└── go.mod
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.20+
- MySQL 8.0
- Docker (optional)

### 1. Clone & configure
```bash
git clone https://github.com/royanqodri/Absensi-APPs.git
cd Absensi-APPs
```

Create a `local.env` file in the project root:
```env
```

### 2. Run locally
```bash
go mod tidy
go run cmd/main.go
```
The app auto-connects to MySQL and runs `AutoMigrate` on startup — no manual migration step needed.

### 3. Run with Docker
```bash
docker build -t absensi-api .
docker run -d -p 8080:8080 --env-file local.env absensi-api
```

### 4. Verify
```bash
curl http://localhost:8080/health
# {"status":"ok","db":"connected"}
```

---

## 📚 API Documentation

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/user` | List employees, with filtering (`created_at`) & pagination (`page_now`, `limit`) |
| POST | `/api/v1/user` | Create or update one or more employees (bulk) |
| GET | `/api/v1/absensi` | List attendance records, filterable by user, date range |
| POST | `/api/v1/absensi` | Create or update one or more attendance records (bulk) |
| GET | `/health` | Service + DB health check |

### Example — Create attendance record
```http

POST /api/v1/absensi
Content-Type: application/json

-- Insert for example :

{
  "data": [
    {
      "id_user": 1,
      "overtime_masuk": "30",
      "overtime_pulang": "45",
      "jam_masuk": "08:00:00",
      "jam_keluar": "17:00:00"
    }
  ]
}

-- update for example :

{
  "data": [
    {
      "id": 1,
      "id_user": 1,
      "overtime_masuk": "30",
      "overtime_pulang": "45",
      "jam_masuk": "08:00:00",
      "jam_keluar": "17:00:00"
    }
  ]
}
```

**Success POST response:**
```json
{
    "status_response": {
        "status_code": 200,
        "status_text": "OK",
        "message": "success",
        "errors": []
    },
    "total_page": 0,
    "total_data": 0,
    "data": []
}
```


**Duplicate error POST response:**
```json
{
    "status_response": {
        "status_code": 409,
        "status_text": "Bad Request",
        "message": "Conflict: Duplicate entry detected",
        "errors": []
    },
    "total_page": 0,
    "total_data": 0,
    "data": []
}
```

**Internal error POST response:**
```json
{
    "status_response": {
        "status_code": 500,
        "status_text": "Internal Server Erro",
        "message": "site not nul",
        "errors": []
    },
    "total_page": 0,
    "total_data": 0,
    "data": []
}
``` 


GET /api/v1/absensi
Content-Type: application/json

**Success GET response:**
```json
{
    "status_response": {
        "status_code": 200,
        "status_text": "OK",
        "message": "success",
        "errors": []
    },
    "total_page": 1,
    "total_data": 2,
    "data": [
        {
            "id": 1,
            "id_user": 1,
            "name": "your_name",
            "username": "username",
            "overtime_masuk": "30",
            "overtime_pulang": "45",
            "jam_masuk": "08:00:00",
            "jam_keluar": "17:00:00",
            "created_at": "2026-08-05T22:17:33.679+07:00",
            "updated_at": "2026-08-05T22:33:33.481+07:00"
        }
    ]
}
```

---

## 🐳 Deployment

The repo ships with a GitHub Actions workflow (`.github/workflows/deploy.yml`) that on every push to `main`:
1. Builds the Docker image
2. Pushes it to Docker Hub
3. SSHes into the target server, pulls the new image, and restarts the container

```bash
# Manual build
docker build -t absensi-api:latest .
docker run -d -p 8080:8080 --env-file local.env absensi-api:latest
```

---

## 🗺️ Roadmap

- [ ] Wire up `/auth/login` and `/auth/register`, and attach `JWTMiddleware()` to protected routes
- [ ] Fix and activate Redis client initialization for caching
- [ ] Add rate limiting
- [ ] Add unit tests (target 70%+ coverage) and a test step in CI
- [ ] Expose generated Swagger/OpenAPI docs (godoc annotations already in controllers)
- [ ] Role-based access control (admin/employee)

---

## 👨‍💻 Author

**Royan Qodri**
[GitHub](https://github.com/royanqodri) · Open to freelance backend work (Go, REST APIs, MySQL)

---

*Built with Go.*
