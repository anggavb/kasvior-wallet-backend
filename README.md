# Kasvior Wallet Backend

![Go](https://img.shields.io/badge/Go-1.26.3-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12.0-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)

<p align="center">
  <img src="https://raw.githubusercontent.com/anggavb/kasvior-wallet-app/refs/heads/main/public/money-wallet-black.png#gh-light-mode-only" alt="Kasvior Logo" />
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/anggavb/kasvior-wallet-app/refs/heads/main/public/money-wallet-black.png#gh-dark-mode-only" alt="Kasvior Logo" />
</p>

Kasvior Wallet Backend is a REST API service for a digital wallet application. It handles authentication, user profile management, wallet dashboard data, top up transactions, transfers, transaction history, payment methods, and API documentation through Swagger.

## Technology Stack

![Go](https://img.shields.io/badge/Go-1.26.3-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=flat-square&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=flat-square&logo=redis&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat-square&logo=jsonwebtokens&logoColor=white)
![SMTP](https://img.shields.io/badge/SMTP-Email_Service-FF9900?style=flat-square&logo=mailchimp&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-API_Docs-85EA2D?style=flat-square&logo=swagger&logoColor=black)
![Docker](https://img.shields.io/badge/Docker-Containerization-2496ED?style=flat-square&logo=docker&logoColor=white)
![golang-migrate](https://img.shields.io/badge/golang--migrate-Database_Migrations-FF6C37?style=flat-square&logo=golang&logoColor=white)

## Features

- User registration, login, logout, forgot password, and reset password
- JWT-based protected routes with active token validation
- User profile, password, and PIN management
- Wallet balance, income, expense, and transaction report endpoints
- Top up transaction flow with payment methods
- Wallet transfer flow with receiver search
- Transaction history with pagination support
- Static file serving for profile images and payment logos
- Swagger documentation available at `/swagger/index.html`
- Database migrations and seeders through Makefile commands

## Prerequisites

Make sure these tools and services are available before running the project locally:

- Go 1.26.3 or later
- PostgreSQL
- Redis
- Make
- `golang-migrate`
- `psql` or any PostgreSQL client like DBeaver, TablePlus, or pgAdmin
- Git

## Setup Instruction

Clone the repository:

```bash
git clone https://github.com/anggavb/kasvior-wallet-backend.git
cd kasvior-wallet-backend
```

Install Go dependencies:

```bash
go mod download
```

Create a `.env` file in the project root:

```env
APP_HOST=localhost
APP_PORT=8080

DB_URL=postgres://postgres:password@localhost:5432/kasvior_wallet?sslmode=disable

RDB_ADDR=localhost:6379
RDB_USER=
RDB_PASS=
RDB_PREFIX=kasvior:

JWT_SECRET=your_jwt_secret
JWT_ISSUER=kasvior-wallet

RESET_PASSWORD_URL=http://localhost:5173/reset-password

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=no-reply@example.com
```

Run database migrations:

```bash
make migrate-up
```

Seed payment methods:

```bash
make seed
```

Run the application locally:

```bash
make run
```

The API will run at:

```text
http://localhost:8080
```

Open Swagger documentation at:

```text
http://localhost:8080/swagger/index.html
```

## Project Structure

This project follows a MVC (Model-View-Controller) pattern, separating concerns into different layers. The main directories are:
```text
.
├── cmd/                 # Application entry point
├── db/
│   ├── migrations/      # SQL migration files
│   └── seeds/           # SQL seeder files
├── docs/                # Generated Swagger documentation
├── internal/
│   ├── apperrors/       # Application error helpers
│   ├── binder/          # Request binding and validation helpers
│   ├── config/          # PostgreSQL and Redis configuration
│   ├── controller/      # HTTP controllers
│   ├── dto/             # Request and response DTOs
│   ├── jwttoken/        # JWT utility helpers
│   ├── middleware/      # HTTP middleware
│   ├── model/           # Domain models
│   ├── repository/      # Database and cache repositories
│   ├── response/        # Response helpers
│   ├── router/          # Route definitions
│   └── service/         # Business logic
├── pkg/                 # Shared packages
├── public/              # Static images and payment logos
├── Dockerfile           # Docker build configuration
├── Makefile             # Development commands
├── go.mod
└── go.sum
```

## How to Contribute

1. Fork this repository.
2. Create a new branch from `main`.
3. Make your changes with clear commit messages.
4. Run the project locally and verify the affected API behavior.
5. Open a pull request with a concise description of the changes.

## Related Project

- [Kasvior Wallet Frontend](https://github.com/anggavb/kasvior-wallet-app)

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
