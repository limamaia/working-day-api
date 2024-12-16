# Task Management API

## Description
This is a RESTful API for managing tasks and users, developed using Golang with the Gin framework. The project follows the MVC architecture with Clean Code principles and includes unit tests.

The API supports JWT authentication, CRUD operations for tasks, and basic user management. It also provides an interactive API documentation via Swagger.

## Requirements
- Docker and Docker Compose

## How to Run

### 1. Run the Application
To build and run the application using Docker Compose, execute:

```bash
docker-compose up --build
```

This command will build the image and start the application along with its dependencies.

## How to Test
To execute unit tests for the services, run:

```bash
go test ./internal/services -v
```

# API Documentation

The Swagger interactive documentation is available at:  
[http://localhost/swagger/index.html#/](http://localhost/swagger/index.html#/)

## Project Structure
```bash
├── api
│   └── swagger            # Swagger documentation files
│       ├── auth_docs.go   # Swagger docs for Auth routes
│       ├── task_docs.go   # Swagger docs for Task routes
│       └── user_docs.go   # Swagger docs for User routes
├── cmd
│   └── main.go            # Entry point of the application
├── config                 # Application configuration
├── database               # Database connection
├── docker                 # Docker-related files
├── docs                   # Generated Swagger docs (swagger.json, swagger.yaml)
├── internal
│   ├── container          # Dependency injection
│   ├── controllers        # API Handlers (Controllers)
│   ├── domain             # Structs and models
│   ├── helpers            # Utility functions (e.g., JWT, password hashing)
│   ├── middlewares        # Middleware (e.g., JWT Authentication)
│   ├── repositories       # Database interaction logic
│   └── services           # Business logic (services) and tests
│       ├── task_service.go
│       ├── user_service.go
│       ├── login_service.go
│       ├── task_service_test.go
│       ├── user_service_test.go
│       └── login_service_test.go
├── messaging              # RabbitMQ integration
├── router                 # Application routes
├── vendor                 # Vendor dependencies
├── .env                   # Environment variables
├── docker-compose.yml     # Docker Compose configuration
├── go.mod                 # Go modules
├── go.sum                 # Go dependencies
└── README.md              # Project documentation
```
The project is organized with **MVC (Model-View-Controller)** architecture and follows **clean code practices**:

## Key Features

- **JWT Authentication**: Secure access to endpoints.
- **Task Management**: Perform CRUD operations for tasks.
- **User Management**: Create, retrieve, and validate users.
- **Swagger Documentation**: Interactive API documentation.
- **Unit Testing**: Comprehensive test coverage for services.