# 🧠 Thera API

**Thera API** is the backend service powering Thera Data Solution’s applications.  
It is built with **Go (Golang)**, using **GORM** as the ORM layer and **PostgreSQL** as the primary database.

This project provides a fast, secure, and scalable RESTful API foundation for applications such as **ERP**, **CRM**, **Booking Systems**, and **E-commerce Platforms**.

---

## ⚙️ Tech Stack

- **Language:** [Go](https://golang.org/)  
- **Framework / ORM:** [GORM](https://gorm.io/)  
- **Database:** [PostgreSQL](https://www.postgresql.org/)  
- **Authentication:** JWT-based middleware  
- **Containerization:** Docker-ready  
- **Architecture:** Clean architecture pattern with modular structure  

---

## 📂 Project Structure

```
thera-api/
├── cmd/                # Application entry points
├── config/             # Configuration and environment setup
├── controllers/        # HTTP handlers for each module
├── middleware/         # JWT, logging, and error handling
├── models/             # GORM models and database relations
├── routes/             # API routing definitions
├── services/           # Business logic layer
├── utils/              # Helpers and utilities
└── main.go             # Application bootstrap
```

---

## 🚀 Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/thera-data/thera-api.git
cd thera-api
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory:
```bash
APP_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/thera_db?sslmode=disable
JWT_SECRET=your-secret-key
```

### 3. Run with Go
```bash
go run main.go
```

### 4. Or Run with Docker
```bash
docker-compose up --build
```

---

## 🧩 Features

- 🔐 Secure authentication with JWT  
- 📦 Modular architecture for scalability  
- ⚡ Fast and lightweight performance  
- 🧾 Auto-migrations with GORM  
- 🌐 Ready for deployment via Docker  

---

## 🧠 Example Endpoints

| Method | Endpoint           | Description              |
|--------|--------------------|--------------------------|
| `POST` | `/api/v1/auth/login`  | User login               |
| `GET`  | `/api/v1/users`       | Get list of users        |
| `POST` | `/api/v1/users`       | Create a new user        |
| `GET`  | `/api/v1/health`      | Health check endpoint    |

---

## 🧰 Development Notes

- Follow Go’s standard project layout  
- Use `go mod tidy` to sync dependencies  
- Use `air` for live reload during local development (optional)

---

## 🧭 Roadmap

- [ ] Add role-based access control (RBAC)  
- [ ] Implement GraphQL gateway  
- [ ] Integrate background jobs (e.g., email notifications)  
- [ ] Unit tests & CI/CD pipeline  

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

> _Built with Go for speed, reliability, and simplicity._
