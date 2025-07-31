package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indalyadav56/gogen/internal/template"
)

// File represents a file to be created
type File struct {
	Path         string
	Package      string
	TemplateName string
}

// FileGenerator handles file generation
type FileGenerator struct {
	renderer    *template.Renderer
	projectRoot string
	moduleName  string
	isMonolith  bool
	useGin      bool
	entities    []string
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(renderer *template.Renderer, projectRoot, moduleName string, isMonolith, useGin bool, entities []string) *FileGenerator {
	return &FileGenerator{
		renderer:    renderer,
		projectRoot: projectRoot,
		moduleName:  moduleName,
		isMonolith:  isMonolith,
		useGin:      useGin,
		entities:    entities,
	}
}

// GenerateFiles creates all project files for the given entity
func (fg *FileGenerator) GenerateFiles(entityName string) error {
	files := fg.getFileList(entityName)
	
	for _, file := range files {
		if err := fg.createFile(file, entityName); err != nil {
			return fmt.Errorf("failed to create file %s: %w", file.Path, err)
		}
	}
	
	return nil
}

// getFileList returns the list of files to create
func (fg *FileGenerator) getFileList(entityName string) []File {
	if fg.isMonolith {
		return fg.getMonolithFileList(entityName)
	}
	return fg.getMicroserviceFileList(entityName)
}

// getMicroserviceFileList returns files for microservice architecture
func (fg *FileGenerator) getMicroserviceFileList(entityName string) []File {
	return []File{
		{Path: "cmd/main.go", Package: "main", TemplateName: "main.tmpl"},
		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},

		// Domain
		{Path: "internal/domain/constants/constants.go", Package: "constants", TemplateName: ""},
		{Path: "internal/domain/entity/entity.go", Package: "entity", TemplateName: "entity.tmpl"},
		{Path: "internal/domain/repository/repository.go", Package: "repository", TemplateName: "repository.tmpl"},

		{Path: "internal/application/" + fmt.Sprintf("%s_service.go", strings.ToLower(entityName)), Package: "application", TemplateName: "service.tmpl"},
		{Path: "internal/interface/http/v1/routes/routes.go", Package: "routes", TemplateName: "routes.tmpl"},
		{Path: "internal/interface/http/v1/handlers/" + fmt.Sprintf("%s_handler.go", strings.ToLower(entityName)), Package: "handlers", TemplateName: "handler.tmpl"},
		{Path: "internal/interface/http/middlewares/auth_middleware.go", Package: "middlewares"},

		{Path: "internal/infrastructure/postgres/postgres.go", Package: "postgres", TemplateName: "postgres_repository.tmpl"},

		// DTO
		{Path: "internal/interface/http/v1/dto/request.go", Package: "dto", TemplateName: ""},
		{Path: "internal/interface/http/v1/dto/response.go", Package: "dto", TemplateName: ""},

		{Path: "pkg/db/db.go", Package: "db", TemplateName: "db.tmpl"},

		{Path: ".gitignore", Package: "", TemplateName: ""},
		{Path: "Dockerfile", Package: "", TemplateName: "docker.tmpl"},
		{Path: "Taskfile.yaml", Package: "", TemplateName: "taskfile.tmpl"},
	}
}

// getMonolithFileList returns files for monolith bounded context architecture
func (fg *FileGenerator) getMonolithFileList(entityName string) []File {
	// Base files
	files := []File{
		{Path: "cmd/main.go", Package: "main", TemplateName: fg.getMainTemplate()},
		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},
		{Path: "pkg/db/db.go", Package: "db", TemplateName: "db.tmpl"},
		
		// Shared components
		{Path: "internal/shared/middleware/auth.go", Package: "middleware", TemplateName: ""},
		{Path: "internal/shared/dto/common.go", Package: "dto", TemplateName: ""},
		{Path: "internal/shared/utils/utils.go", Package: "utils", TemplateName: ""},
		
		{Path: ".gitignore", Package: "", TemplateName: ""},
		{Path: "Dockerfile", Package: "", TemplateName: "docker.tmpl"},

		// DTO
		{Path: "internal/"+entityName+"/interface/http/v1/dto/request.go", Package: "dto", TemplateName: ""},
		{Path: "internal/"+entityName+"/interface/http/v1/dto/response.go", Package: "dto", TemplateName: ""},
	}
	
	// If no entity specified, create example bounded context
	if entityName == "" {
		entityName = "example"
	}
	
	// Create bounded context files for the entity
	entityLower := strings.ToLower(entityName)
	entityPath := fmt.Sprintf("internal/%s", entityLower)
	
	boundedContextFiles := []File{
		// Domain layer - core business entities, value objects, aggregates
		{Path: entityPath + "/domain/entity/entity.go", Package: "entity", TemplateName: "entity.tmpl"},
		{Path: entityPath + "/domain/repository/repository.go", Package: "repository", TemplateName: "repository.tmpl"},
		
		// Interface layer - HTTP handlers, controllers (dynamic based on framework)
		{Path: entityPath + "/interface/http/v1/handlers/" + fmt.Sprintf("%s_handler.go", strings.ToLower(entityName)), Package: "handlers", TemplateName: fg.getHandlerTemplate()},
		{Path: entityPath + "/interface/http/v1/routes/routes.go", Package: "routes", TemplateName: fg.getRoutesTemplate()},
		
		// Application layer - application services, use cases
		{Path: entityPath + "/application/" + fmt.Sprintf("%s_service.go", strings.ToLower(entityName)), Package: "application", TemplateName: "service.tmpl"},
		
		// Infrastructure layer - database implementations, external APIs
		{Path: entityPath + "/infrastructure/postgres/postgres.go", Package: "postgres", TemplateName: "postgres_repository.tmpl"},
		
	}
	
	return append(files, boundedContextFiles...)
}

// getHandlerTemplate returns the appropriate handler template based on framework
func (fg *FileGenerator) getHandlerTemplate() string {
	if fg.useGin {
		return "gin_handler.tmpl"
	}
	return "handler.tmpl"
}

// getRoutesTemplate returns the appropriate routes template based on framework
func (fg *FileGenerator) getRoutesTemplate() string {
	if fg.useGin {
		return "gin_routes.tmpl"
	}
	return "routes.tmpl"
}

// getMainTemplate returns the appropriate main template based on framework and architecture
func (fg *FileGenerator) getMainTemplate() string {
	if fg.isMonolith {
		if fg.useGin {
			return "gin_monolith_main.tmpl"
		}
		return "monolith_main.tmpl"
	}
	return "main.tmpl"
}

// prepareTemplateData creates template data with correct import paths based on architecture
func (fg *FileGenerator) prepareTemplateData(packageName, entityName string) template.Data {
	data := template.Data{
		Package:     packageName,
		ProjectRoot: fg.projectRoot,
		ModuleName:  fg.moduleName,
		EntityName:  entityName,
		Entities:    fg.entities,
		IsMonolith:  fg.isMonolith,
		UseGin:      fg.useGin,
	}
	
	if fg.isMonolith && entityName != "" {
		// Monolith bounded context import paths (clean architecture)
		entityLower := strings.ToLower(entityName)
		data.HandlerImport = fmt.Sprintf("%s/internal/%s/interface/http/v1/handlers", fg.moduleName, entityLower)
		data.ServiceImport = fmt.Sprintf("%s/internal/%s/application", fg.moduleName, entityLower)
		data.RepositoryImport = fmt.Sprintf("%s/internal/%s/domain/repository", fg.moduleName, entityLower)
		data.EntityImport = fmt.Sprintf("%s/internal/%s/domain/entity", fg.moduleName, entityLower)
		data.InfraImport = fmt.Sprintf("%s/internal/%s/infrastructure", fg.moduleName, entityLower)
		data.RoutesImport = fmt.Sprintf("%s/internal/%s/interface/http/v1/routes", fg.moduleName, entityLower)
	} else {
		// Microservice import paths (default)
		data.HandlerImport = fmt.Sprintf("%s/internal/interface/http/v1/handlers", fg.moduleName)
		data.ServiceImport = fmt.Sprintf("%s/internal/application", fg.moduleName)
		data.RepositoryImport = fmt.Sprintf("%s/internal/domain/repository", fg.moduleName)
		data.EntityImport = fmt.Sprintf("%s/internal/domain/entity", fg.moduleName)
		data.InfraImport = fmt.Sprintf("%s/internal/infrastructure/postgres", fg.moduleName)
		data.RoutesImport = fmt.Sprintf("%s/internal/interface/http/v1/routes", fg.moduleName)
	}
	
	return data
}

// createFile creates a single file
func (fg *FileGenerator) createFile(file File, entityName string) error {
	fullPath := filepath.Join(fg.projectRoot, file.Path)
	
	// Prepare template data
	templateData := fg.prepareTemplateData(file.Package, entityName)
	
	// If a template is specified, render it
	if file.TemplateName != "" && entityName != "" {
		templatePath := "templates/" + file.TemplateName
		return fg.renderer.RenderToFile(templatePath, fullPath, templateData)
	} else if file.TemplateName == "db.tmpl" || file.TemplateName == "logger.tmpl" {
		templatePath := "templates/" + file.TemplateName
		return fg.renderer.RenderToFile(templatePath, fullPath, templateData)
	} else {
		content := fmt.Sprintf("package %s\n", file.Package)
		if file.Package == "" {
			return os.WriteFile(fullPath, nil, 0644)
		}
		return os.WriteFile(fullPath, []byte(content), 0644)
	}
}
