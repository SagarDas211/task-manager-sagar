# Task Management API (Golang)

A cleanly structured Task Management REST API built using **Golang**, **Gin**, and **Domain-Driven Design (DDD)** principles.
This project demonstrates **layered architecture**, **test-driven development (TDD)**, and **clean coding practices**.

---

## 🚀 Features

* Create, update, delete, and retrieve tasks
* In-memory data storage (no DB required)
* Proper layering: Domain, Repository, Service, Handler
* Input validation and error handling
* Filtering and pagination support
* Unit tests (service layer)
* Integration tests (HTTP endpoints)

---

## 🧱 Project Structure

```
task-manager/
│
├── cmd/
│   └── main.go                # Application entry point
│
├── internal/
│   ├── domain/               # Core business entities
│   ├── service/              # Business logic layer
│   ├── handler/              # HTTP handlers (Gin)
│   ├── repository/           # Data persistence (in-memory)
│   ├── errors/               # Custom error definitions
│
├── tests/                    # (optional) extra tests (currently empty)
├── go.mod
└── README.md
```

---

## ⚙️ Tech Stack

* Go (Golang, Gin)
* Gin (HTTP framework)
* Standard library (`net/http`, `testing`, etc.)

---

## 🛠️ Setup Instructions

### 1. Clone the Repository

```
git clone https://github.com/SagarDas211/task-manager-sagar.git
cd task-manager-sagar
```

---

### 2. Install Dependencies

```
go mod tidy
```

---

### 3. Run the Application

```
go run ./cmd
```

Server will start at:

```
http://localhost:8080
```

---

## 🧪 Running Tests

### Run all tests:

```
go test ./...
```

### Run with verbose output:

```
go test ./... -v
```

### Run only service package tests:

```
go test ./internal/service -v
```

### Run only handler package tests:

```
go test ./internal/handler -v
```

### Run a single test by name:

```
go test ./internal/service -run TestTaskService_CreateTask_Success -v
```

If you see a cache permission error in some environments, run:

```
GOCACHE=/tmp/go-build go test ./...
```

---

## 📡 API Endpoints

---

### ➕ Create Task

**POST** `/tasks`

#### Request

```json
{
  "title": "Complete assignment",
  "description": "Finish DDD project",
  "due_date": "2030-01-01"
}
```

#### Response (201)

```json
{
  "id": "uuid",
  "title": "Complete assignment",
  "description": "Finish DDD project",
  "status": "PENDING",
  "due_date": "2030-01-01"
}
```

---

### 🔍 Get Task

**GET** `/tasks/{id}`

#### Response (200)

```json
{
  "id": "uuid",
  "title": "...",
  "status": "PENDING"
}
```

#### Not Found (404)

```json
{
  "error": "task not found"
}
```

---

### ✏️ Update Task

**PUT** `/tasks/{id}`

#### Request

```json
{
  "title": "Updated Title",
  "status": "IN_PROGRESS"
}
```

#### Response (200)

```json
{
  "id": "...",
  "title": "Updated Title"
}
```

---

### ❌ Delete Task

**DELETE** `/tasks/{id}`

#### Response

```
204 No Content
```

---

### 📋 List Tasks

**GET** `/tasks`

#### Query Parameters

| Param  | Description       |
| ------ | ----------------- |
| status | Filter by status  |
| limit  | Number of results |
| offset | Pagination offset |

Tasks are always sorted by `due_date`.

---

#### Example

```
GET /tasks?status=PENDING&limit=5&offset=0
```

#### Response (200)

```json
[
  {
    "id": "...",
    "title": "...",
    "status": "PENDING"
  }
]
```

---

## 🧠 Design Decisions

* **DDD Layers**

  * Domain → core business rules
  * Service → business logic
  * Repository → persistence abstraction
  * Handler → HTTP layer

* **Validation Split**

  * Handler → request format validation
  * Service/Domain → business validation

* **Thread-Safe Repository**

  * Uses `sync.RWMutex` for concurrent access

* **Pagination & Filtering**

  * Implemented in service layer (not repository)

---

## ⚠️ Assumptions

* In-memory storage (data resets on restart)
* No authentication required
* UUID used for ID generation

---

## 🧪 Testing Strategy

* **Unit Tests**

  * Service layer with mocked repository

* **Integration Tests**

  * Full HTTP testing using Gin + httptest

---

## 🚀 Future Improvements

* Persistent database (PostgreSQL)
* Authentication & authorization
* Docker support
* OpenAPI/Swagger documentation
* Cursor-based pagination

---

## 👨‍💻 Author

Sagar Das

---

## 📄 License

This project is for evaluation/demo purposes.
