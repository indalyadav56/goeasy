package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds data for file templates
type TemplateData struct {
	ModuleName       string // e.g., auth-service
	PascalModuleName string // e.g., AuthService
}

// templates maps file paths to their content
var templates = map[string]string{
	"cmd/api/main.go": `package main

import (
	"log"
	"{{.ModuleName}}/internal/interface/http"
)

func main() {
	log.Println("Starting {{.PascalModuleName}}...")
	http.StartServer()
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
	"internal/domain/service/service.go": `package service

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
	"internal/application/app.go": `package application

import (
	"context"
	"{{.ModuleName}}/internal/domain/service"
)

type {{.PascalModuleName}}App struct {
	svc *service.{{.PascalModuleName}}Service
}

func New{{.PascalModuleName}}App(svc *service.{{.PascalModuleName}}Service) *{{.PascalModuleName}}App {
	return &{{.PascalModuleName}}App{svc: svc}
}

func (a *{{.PascalModuleName}}App) Create{{.PascalModuleName}}(ctx context.Context, name string) error {
	_, err := a.svc.Create(ctx, name)
	return err
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
	"internal/interface/http/handler.go": `package http

import (
	"net/http"
	"{{.ModuleName}}/internal/application"
)

type Handler struct {
	app *application.{{.PascalModuleName}}App
}

func NewHandler(app *application.{{.PascalModuleName}}App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) StartServer() {
	http.HandleFunc("/create", h.createHandler)
	http.ListenAndServe(":8080", nil)
}

func (h *Handler) createHandler(w http.ResponseWriter, r *http.Request) {
	// Implement handler logic
}
`,
	"config/config.yaml": `server:
  port: 8080
database:
  host: localhost
  port: 5432
`,
	"README.md": `# {{.PascalModuleName}}

A microservice built with Domain-Driven Design (DDD) principles.

## Structure
- /cmd: Entry points
- /internal/domain: Entities, repositories
- /internal/application: Use cases
- /internal/application/services: Use cases
- /internal/infrastructure: Database, logging
- /internal/interface: HTTP handlers
- /config: Configuration files

## Setup
1. Install dependencies: ` + "`go mod tidy`" + `
2. Run: ` + "`go run cmd/api/main.go`" + `
`,
	"go.mod": `module {{.ModuleName}}

go 1.21
`,
}

func main() {
	// Parse CLI arguments
	moduleName := flag.String("generate", "", "Module name (e.g., auth-service)")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("Error: Provide a module name using -generate")
		fmt.Println("Example: gogen -generate auth-service")
		os.Exit(1)
	}

	// Validate module name (kebab-case)
	if !isValidKebabCase(*moduleName) {
		fmt.Println("Error: Module name must be in kebab-case (e.g., auth-service)")
		os.Exit(1)
	}

	// Convert to PascalCase
	pascalModuleName := toPascalCase(*moduleName)

	// Create structure
	if err := createModuleStructure(*moduleName, pascalModuleName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated DDD module: %s\n", *moduleName)
}

// isValidKebabCase checks if the name is in kebab-case
func isValidKebabCase(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || c == '-' || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return strings.Contains(name, "-")
}

// toPascalCase converts kebab-case to PascalCase
func toPascalCase(name string) string {
	words := strings.Split(name, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, "")
}

// createModuleStructure creates folders and files
func createModuleStructure(moduleName, pascalModuleName string) error {
	// Create root directory
	if err := os.Mkdir(moduleName, 0755); err != nil {
		return err
	}

	// Template data
	data := TemplateData{
		ModuleName:       moduleName,
		PascalModuleName: pascalModuleName,
	}

	// Create files
	for path, content := range templates {
		fullPath := filepath.Join(moduleName, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer file.Close()

		tmpl, err := template.New(path).Parse(content)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}
	}

	return nil
}
