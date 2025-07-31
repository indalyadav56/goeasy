# 🚀 GoGen - Go Project Generator

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

A powerful CLI tool that generates **production-ready Go projects** with **Clean Architecture**, **Domain-Driven Design (DDD)**, and **best practices**. Supports both **microservice** and **monolith** architectures with **Chi** and **Gin** framework options.

## ✨ Features

- 🏗️ **Clean Architecture**: Perfect 4-layer separation (Domain, Interface, Application, Infrastructure)
- 🎯 **Domain-Driven Design**: Bounded contexts with proper entity isolation
- 🔄 **Multiple Architectures**: Microservice and Monolith support
- 🌐 **Framework Choice**: Chi Router (default) or Gin framework
- 📦 **Multi-Entity Support**: Generate multiple bounded contexts
- 🐳 **Docker Ready**: Includes optimized Dockerfile
- 📝 **Template System**: Extensible and customizable templates
- 🔧 **Go Modules**: Automatic module initialization and dependency management

## 🚀 Quick Start

### Installation

```bash
# Install from source
go install github.com/indalyadav56/gogen/cmd/gogen@latest

# Or build locally
git clone https://github.com/indalyadav56/gogen.git
cd gogen
go build -o gogen cmd/gogen/main.go
```

### Basic Usage

```bash
# Generate microservice with Chi router
gogen --module github.com/yourname/project --entity user

# Generate monolith with multiple entities (Chi)
gogen --module github.com/yourname/project --entity user --entity product --monolith

# Generate monolith with Gin framework
gogen --module github.com/yourname/project --entity user --entity order --monolith --gin
```

## 📖 Usage Examples

### Microservice Architecture

Generate a single-entity microservice:

```bash
gogen --module github.com/company/user-service --entity user
```

**Generated Structure:**
```
user-service/
├── cmd/main.go
├── internal/
│   ├── application/user_service.go
│   ├── domain/
│   │   ├── entity/entity.go
│   │   └── repository/repository.go
│   ├── interface/http/v1/
│   │   ├── handlers/user_handler.go
│   │   └── routes/routes.go
│   └── infrastructure/postgres/postgres.go
├── pkg/
│   ├── db/db.go
│   └── logger/logger.go
├── config/config.go
├── Dockerfile
└── go.mod
```

### Monolith Architecture with Clean Architecture

Generate a monolith with multiple bounded contexts:

```bash
gogen --module github.com/company/ecommerce --entity user --entity product --entity order --monolith
```

**Generated Structure:**
```
ecommerce/
├── cmd/main.go
├── internal/
│   ├── shared/          # Cross-cutting concerns
│   │   ├── middleware/
│   │   ├── dto/
│   │   └── utils/
│   ├── user/            # User bounded context
│   │   ├── domain/
│   │   │   ├── entity/entity.go
│   │   │   └── repository/repository.go
│   │   ├── interface/http/v1/
│   │   │   ├── handlers/user_handler.go
│   │   │   ├── routes/routes.go
│   │   │   └── dto/
│   │   ├── application/user_service.go
│   │   └── infrastructure/postgres/postgres.go
│   ├── product/         # Product bounded context
│   │   ├── domain/
│   │   ├── interface/
│   │   ├── application/
│   │   └── infrastructure/
│   └── order/           # Order bounded context
│       ├── domain/
│       ├── interface/
│       ├── application/
│       └── infrastructure/
├── pkg/
├── config/
├── migrations/
├── scripts/
└── Dockerfile
```

### Gin Framework Support

Generate with Gin instead of Chi:

```bash
gogen --module github.com/company/api --entity user --monolith --gin
```

This generates Gin-specific handlers, routes, and main.go with Gin router setup.

## 🛠️ Command Line Options

| Flag | Description | Example |
|------|-------------|----------|
| `--module` | Go module name | `github.com/user/project` |
| `--entity` | Entity name (can be used multiple times) | `--entity user --entity product` |
| `--monolith` | Generate monolith architecture | `--monolith` |
| `--gin` | Use Gin framework instead of Chi | `--gin` |

## 🏗️ Architecture Patterns

### Clean Architecture Layers

1. **Domain Layer** (`domain/`)
   - Entities, Value Objects, Aggregates
   - Repository interfaces
   - Business rules and logic

2. **Interface Layer** (`interface/`)
   - HTTP handlers and controllers
   - Route definitions
   - DTOs and request/response models
   - Middleware

3. **Application Layer** (`application/`)
   - Use cases and application services
   - Orchestrates domain objects
   - Implements business workflows

4. **Infrastructure Layer** (`infrastructure/`)
   - Database implementations
   - External API clients
   - Framework-specific code

### Framework Support

#### Chi Router (Default)
- Lightweight and fast
- Composable middleware
- RESTful routing

#### Gin Framework (`--gin` flag)
- High performance
- Built-in middleware
- JSON binding and validation

## 📁 Generated Files

### Core Files
- `cmd/main.go` - Application entry point
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `.gitignore` - Git ignore patterns

### Business Logic
- `*_service.go` - Application services
- `*_handler.go` - HTTP handlers
- `entity.go` - Domain entities
- `repository.go` - Repository interfaces
- `postgres.go` - Database implementations

### Infrastructure
- `config/config.go` - Configuration management
- `pkg/db/db.go` - Database connection
- `pkg/logger/logger.go` - Logging setup
- `migrations/` - Database migrations

## 🔧 Development

### Prerequisites
- Go 1.21 or higher
- Git

### Building from Source

```bash
git clone https://github.com/indalyadav56/gogen.git
cd gogen
go mod tidy
go build -o gogen cmd/gogen/main.go
```

### Running Tests

```bash
go test ./...
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Clean Architecture principles by Robert C. Martin
- Domain-Driven Design concepts by Eric Evans
- Go community best practices

## 📞 Support

If you have any questions or need help, please:
- Open an issue on GitHub
- Check the documentation
- Join our community discussions

---

**Made with ❤️ for the Go community**
