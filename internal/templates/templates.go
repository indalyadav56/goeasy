package templates

// Templates maps file paths to their content
var Templates = map[string]string{
	"cmd/api/main.go": `package main

import (
	"log"
	"{{.ModuleName}}/internal/interface/http/v1/handlers"
)

func main() {
	log.Println("Starting {{.PascalModuleName}}...")
	handlers.StartServer()
}
`,
	"config/config.yaml": `server:
  port: 8080
database:
  host: localhost
  port: 5432
`,
	"internal/application/services/service.go": `package services

import (
	"context"
	"{{.ModuleName}}/internal/domain/entity"
	"{{.ModuleName}}/internal/domain/repository"
)

type {{.PascalModuleName}}Service struct {
	repo repository.{{.PascalModuleName}}Repository
}

func New{{.PascalModuleName}}Service(repo repository.{{.PascalModuleName}}Repository) *{{.PascalModuleName}}Service {
	return &{{.PascalModuleName}}Service{repo: repo}
}

func (s *{{.PascalModuleName}}Service) Create(ctx context.Context, name string) (*entity.{{.PascalModuleName}}, error) {
	e := &entity.{{.PascalModuleName}}{ID: "generated-id", Name: name}
	return e, s.repo.Save(ctx, e)
}
`,
	"internal/domain/entity/entity.go": `package entity

type {{.PascalModuleName}} struct {
	ID   string
	Name string
}
`,
	"internal/domain/repository/repository.go": `package repository

import (
	"context"
	"{{.ModuleName}}/internal/domain/entity"
)

type {{.PascalModuleName}}Repository interface {
	Save(ctx context.Context, e *entity.{{.PascalModuleName}}) error
	FindByID(ctx context.Context, id string) (*entity.{{.PascalModuleName}}, error)
}
`,
	"internal/infrastructure/db/db.go": `package db

import (
	"context"
	"{{.ModuleName}}/internal/domain/entity"
	"{{.ModuleName}}/internal/domain/repository"
)

type {{.PascalModuleName}}RepositoryImpl struct{}

func New{{.PascalModuleName}}Repository() repository.{{.PascalModuleName}}Repository {
	return &{{.PascalModuleName}}RepositoryImpl{}
}

func (r *{{.PascalModuleName}}RepositoryImpl) Save(ctx context.Context, e *entity.{{.PascalModuleName}}) error {
	// Implement database logic
	return nil
}

func (r *{{.PascalModuleName}}RepositoryImpl) FindByID(ctx context.Context, id string) (*entity.{{.PascalModuleName}}, error) {
	// Implement database logic
	return &entity.{{.PascalModuleName}}{ID: id, Name: "example"}, nil
}
`,
	"internal/infrastructure/logging/logger.go": `package logging

import "log"

func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
`,
	"internal/interface/http/v1/handlers/handler.go": `package handlers

import (
	"net/http"
	"{{.ModuleName}}/internal/application/services"
)

type Handler struct {
	svc *services.{{.PascalModuleName}}Service
}

func NewHandler(svc *services.{{.PascalModuleName}}Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) StartServer() {
	http.HandleFunc("/v1/create", h.createHandler)
	http.ListenAndServe(":8080", nil)
}

func (h *Handler) createHandler(w http.ResponseWriter, r *http.Request) {
	// Implement handler logic
}
`,
	"internal/interface/http/v1/routes/routes.go": `package routes

import (
	"net/http"
	"{{.ModuleName}}/internal/interface/http/v1/handlers"
)

func SetupRoutes(handler *handlers.Handler) {
	http.HandleFunc("/v1/create", handler.createHandler)
}
`,
	"go.mod": `module {{.ModuleName}}

go 1.21
`,
	"README.md": `# {{.PascalModuleName}}

A microservice built with Domain-Driven Design (DDD) principles.

## Structure
- /cmd: Entry points
- /internal/domain: Entities and repositories
- /internal/application/services: Application services
- /internal/infrastructure: Database and logging
- /internal/interface/http/v1: HTTP handlers and routes
- /config: Configuration files

## Setup
1. Install dependencies: ` + "`go mod tidy`" + `
2. Run: ` + "`go run cmd/api/main.go`" + `
`,
}
