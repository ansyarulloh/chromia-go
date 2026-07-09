# Chromia Core Backend API

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-000000?style=for-the-badge&logo=go&logoColor=white)
![Gemini AI](https://img.shields.io/badge/Gemini_AI-8E75B2?style=for-the-badge&logo=googlebard&logoColor=white)

The **Core Backend API** for the Chromia Color Intelligence Ecosystem. This repository houses the highly scalable, microservice-ready backend built with Go (Golang) and the Fiber framework. It acts as the central brain of the ecosystem, orchestrating database transactions, user authentication, and AI prompt engineering via the Google Gemini API.

## Architecture & Technology Stack

The backend strictly adheres to RESTful API principles and MVC architecture to ensure maximum scalability, security, and maintainability.

| Layer | Technologies |
| :--- | :--- |
| **Framework** | Go (Golang) + Fiber (Express-inspired web framework) |
| **ORM & Persistence** | GORM connecting to PostgreSQL 13+ |
| **Security & Auth** | JWT (JSON Web Tokens) with Middleware validation |
| **AI Processing** | Google Generative AI SDK (DeepMind Gemini) |
| **Environment** | Godotenv for secure `.env` secret management |

## Core Features

- **Blazingly Fast Routing:** Built on top of Fiber/Fasthttp, ensuring ultra-low latency responses to the mobile client.
- **AI-Powered Diagnostics (Gemini):** A robust integration with the Gemini AI model. It receives image payloads, constructs context-aware prompts, and extracts precise Color Names and Hex Codes with structured JSON outputs.
- **High-Security Authentication:** JWT-based authentication with role-based access control. All protected routes are guarded by strict middleware validation.
- **Smart Quota Management:** Automated logic to calculate daily scan quotas for users, instantly resetting quotas when a new day begins.
- **Automated Database Migration:** GORM `AutoMigrate` is configured to seamlessly build and update the `users` and `scan_histories` schemas upon server startup.

## Installation & Setup Guide

### System Requirements:

- Go v1.21 or higher
- PostgreSQL v13 or higher (or Docker & Docker Compose)

### 1. Database Setup

We use Docker to spin up the PostgreSQL database instantly. Ensure Docker is running, then execute the following command in the root directory:

```bash
docker-compose up -d
```

This will automatically pull the PostgreSQL image and create a container with the `chromia_db` database, user `chromia_admin`, and password `chromiasatset123` on port 5432.

### 2. Repository Setup

Clone the repository and configure your environment:

```bash
git clone https://github.com/your-username/chromia-go.git
cd chromia-go
```

Create a `.env` file in the root directory and add the following secrets:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=chromia_admin
DB_PASS=chromiasatset123
DB_NAME=chromia_db

# Gemini API Key
GEMINI_API_KEY=your_google_gemini_api_key_here

# JWT Secret
JWT_SECRET=super_secret_chromia_key_2026
```

### 3. Run the Server

Install Go dependencies and spin up the server:

```bash
go mod tidy
go run main.go
```
*Note: The backend runs on `localhost:3000` by default. You will see "Database Migration sukses!" in your terminal if everything is configured correctly.*

## Production Build

To compile the Backend for production deployment (Linux/Ubuntu example):

```bash
GOOS=linux GOARCH=amd64 go build -o chromia-server
./chromia-server
```

## Documentation & Standards

This repository is an Industry-Ready boilerplate. It rigorously follows engineering standards ensuring:

- **Strict Controller/Model Separation:** Business logic is isolated in the `controllers` directory, while database schemas are strictly typed in `models`.
- **Graceful Error Handling:** Standardized JSON error responses (`{"status": "error", "message": "..."}`) are enforced across all endpoints.
- **Ignored Secrets:** All environment secrets and generated binaries are safely ignored in `.gitignore`.
