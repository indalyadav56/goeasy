package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indalyadav56/gogen/internal/template"
)

type File struct {
	Path         string
	Package      string
	TemplateName string
}

type FileGenerator struct {
	renderer    *template.Renderer
	projectRoot string
	moduleName  string
	isMonolith  bool
	useGin      bool
	useAuth     bool
	entities    []string
}

func NewFileGenerator(renderer *template.Renderer, projectRoot, moduleName string, isMonolith, useGin, useAuth bool, entities []string) *FileGenerator {
	return &FileGenerator{
		renderer:    renderer,
		projectRoot: projectRoot,
		moduleName:  moduleName,
		isMonolith:  isMonolith,
		useGin:      useGin,
		useAuth:     useAuth,
		entities:    entities,
	}
}

func (fg *FileGenerator) GenerateFiles(entityName string) error {
	files := fg.getFileList(entityName)
	
	for _, file := range files {
		if err := fg.createFile(file, entityName); err != nil {
			return fmt.Errorf("failed to create file %s: %w", file.Path, err)
		}
	}
	
	return nil
}

func (fg *FileGenerator) getFileList(entityName string) []File {
	if fg.isMonolith {
		return fg.getMonolithFileList(entityName)
	}
	return fg.getMicroserviceFileList(entityName)
}

func (fg *FileGenerator) getMicroserviceFileList(entityName string) []File {
	files := []File{
		{Path: "cmd/main.go", Package: "main", TemplateName: "main.tmpl"},
		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},

		// Domain
		{Path: "internal/domain/constants/constants.go", Package: "constants", TemplateName: ""},
		{Path: "internal/domain/entity/entity.go", Package: "entity", TemplateName: "entity.tmpl"},
		{Path: "internal/domain/repository/repository.go", Package: "repository", TemplateName: "repository.tmpl"},

		{Path: "internal/application/" + fmt.Sprintf("%s_service.go", strings.ToLower(entityName)), Package: "application", TemplateName: "service.tmpl"},
		{Path: "internal/interface/http/v1/routes/routes.go", Package: "routes", TemplateName: fg.getRoutesTemplate()},
		{Path: "internal/interface/http/v1/handlers/" + fmt.Sprintf("%s_handler.go", strings.ToLower(entityName)), Package: "handlers", TemplateName: fg.getHandlerTemplate()},
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

	// Add auth-related files if UseAuth is enabled
	if fg.useAuth {
		authFiles := []File{
			// Auth entities
			{Path: "internal/domain/entity/user.go", Package: "entity", TemplateName: "auth_user_entity.tmpl"},
			{Path: "internal/domain/entity/role.go", Package: "entity", TemplateName: "auth_role_entity.tmpl"},
			{Path: "internal/domain/entity/permission.go", Package: "entity", TemplateName: "auth_permission_entity.tmpl"},
			
			// Auth repositories
			{Path: "internal/domain/repository/user_repository.go", Package: "repository", TemplateName: "auth_user_repository.tmpl"},
			{Path: "internal/domain/repository/role_repository.go", Package: "repository", TemplateName: "auth_role_repository.tmpl"},
			{Path: "internal/domain/repository/permission_repository.go", Package: "repository", TemplateName: "auth_permission_repository.tmpl"},
			
			// Auth services
			{Path: "internal/application/auth_service.go", Package: "application", TemplateName: "auth_service.tmpl"},
			{Path: "internal/application/user_service.go", Package: "application", TemplateName: "auth_user_service.tmpl"},
			{Path: "internal/application/role_service.go", Package: "application", TemplateName: "auth_role_service.tmpl"},
			
			// Auth handlers
			{Path: "internal/interface/http/v1/handlers/auth_handler.go", Package: "handlers", TemplateName: fg.getHandlerTemplate()},
			{Path: "internal/interface/http/v1/handlers/user_handler.go", Package: "handlers", TemplateName: fg.getHandlerTemplate()},
			{Path: "internal/interface/http/v1/handlers/role_handler.go", Package: "handlers", TemplateName: fg.getHandlerTemplate()},
			
			// Auth routes
			{Path: "internal/interface/http/v1/routes/auth_routes.go", Package: "routes", TemplateName: fg.getRoutesTemplate()},
			{Path: "internal/interface/http/v1/routes/role_routes.go", Package: "routes", TemplateName: fg.getRoutesTemplate()},
			
			// Auth infrastructure
			{Path: "internal/infrastructure/postgres/user_postgres.go", Package: "postgres", TemplateName: "auth_user_postgres.tmpl"},
			{Path: "internal/infrastructure/postgres/role_postgres.go", Package: "postgres", TemplateName: "auth_role_postgres.tmpl"},
			{Path: "internal/infrastructure/postgres/permission_postgres.go", Package: "postgres", TemplateName: "auth_permission_postgres.tmpl"},
			
			// Auth middleware and utilities
			{Path: "pkg/auth/jwt.go", Package: "auth", TemplateName: "auth_jwt.tmpl"},
			{Path: "pkg/auth/password.go", Package: "auth", TemplateName: "auth_password.tmpl"},
			{Path: "pkg/auth/rbac.go", Package: "auth", TemplateName: "auth_rbac.tmpl"},
			
			// Auth migrations
			{Path: "migrations/001_create_auth_tables.sql", Package: "", TemplateName: "auth_migration.tmpl"},
		}
		files = append(files, authFiles...)
	}

	return files
}

// getMonolithFileList returns files for monolith bounded context architecture
func (fg *FileGenerator) getMonolithFileList(entityName string) []File {
	// Base files
	files := []File{
		{Path: "cmd/main.go", Package: "main", TemplateName: fg.getMainTemplate()},
		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},
		{Path: "pkg/db/db.go", Package: "db", TemplateName: "db.tmpl"},
		
		// Shared components (non-auth related)
		{Path: "internal/shared/dto/common.go", Package: "dto", TemplateName: ""},
		{Path: "internal/shared/utils/utils.go", Package: "utils", TemplateName: ""},
		{Path: "internal/shared/middleware/middleware.go", Package: "middleware", TemplateName: "auth_middleware.tmpl"},
		
		{Path: ".gitignore", Package: "", TemplateName: ""},
		{Path: "Dockerfile", Package: "", TemplateName: "docker.tmpl"},
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
		
		// DTO
		{Path: entityPath + "/interface/http/v1/dto/request.go", Package: "dto", TemplateName: ""},
		{Path: entityPath + "/interface/http/v1/dto/response.go", Package: "dto", TemplateName: ""},
	}
	
	// Add auth-related bounded contexts if UseAuth is enabled
	if fg.useAuth {
		// JWT utilities in pkg
		authUtilities := []File{
			{Path: "pkg/auth/jwt.go", Package: "auth", TemplateName: "jwt_utils.tmpl"},
			{Path: "migrations/000001_create_auth_tables.up.sql", Package: "", TemplateName: "auth_migration.tmpl"},
		}
		
		// Auth bounded context (authentication/authorization logic)
		authPath := "internal/auth"
		authFiles := []File{
			// Application layer
			{Path: authPath + "/application/auth_service.go", Package: "application", TemplateName: "auth_service.tmpl"},
			
			// Interface layer
			{Path: authPath + "/interface/http/v1/handlers/auth_handler.go", Package: "handlers", TemplateName: "auth_handler.tmpl"},
			{Path: authPath + "/interface/http/v1/routes/auth_routes.go", Package: "routes", TemplateName: "auth_routes.tmpl"},
			{Path: authPath + "/interface/http/v1/dto/auth_request.go", Package: "dto", TemplateName: "auth_request_dto.tmpl"},
			{Path: authPath + "/interface/http/v1/dto/auth_response.go", Package: "dto", TemplateName: "auth_response_dto.tmpl"},
		}
		
		// User bounded context (user management)
		userPath := "internal/user"
		userFiles := []File{
			// Domain layer
			{Path: userPath + "/domain/entity/user.go", Package: "entity", TemplateName: "user_entity.tmpl"},
			{Path: userPath + "/domain/repository/user_repository.go", Package: "repository", TemplateName: "user_repository.tmpl"},
			
			// Application layer
			{Path: userPath + "/application/user_service.go", Package: "application", TemplateName: "user_service.tmpl"},
			
			// Interface layer
			{Path: userPath + "/interface/http/v1/handlers/user_handler.go", Package: "handlers", TemplateName: "user_handler.tmpl"},
			{Path: userPath + "/interface/http/v1/routes/user_routes.go", Package: "routes", TemplateName: "auth_routes.tmpl"},
			{Path: userPath + "/interface/http/v1/dto/user_request.go", Package: "dto", TemplateName: "auth_request_dto.tmpl"},
			{Path: userPath + "/interface/http/v1/dto/user_response.go", Package: "dto", TemplateName: "auth_response_dto.tmpl"},
			
			// Infrastructure layer
			{Path: userPath + "/infrastructure/postgres/user_postgres.go", Package: "postgres", TemplateName: "user_postgres.tmpl"},
		}
		
		// Role bounded context (role management)
		rolePath := "internal/role"
		roleFiles := []File{
			// Domain layer
			{Path: rolePath + "/domain/entity/role.go", Package: "entity", TemplateName: "role_entity.tmpl"},
			{Path: rolePath + "/domain/repository/role_repository.go", Package: "repository", TemplateName: "role_repository.tmpl"},
			
			// Application layer
			{Path: rolePath + "/application/role_service.go", Package: "application", TemplateName: "role_service.tmpl"},
			
			// Interface layer
			{Path: rolePath + "/interface/http/v1/handlers/role_handler.go", Package: "handlers", TemplateName: "role_handler.tmpl"},
			{Path: rolePath + "/interface/http/v1/routes/role_routes.go", Package: "routes", TemplateName: "role_routes.tmpl"},
			{Path: rolePath + "/interface/http/v1/dto/role_request.go", Package: "dto", TemplateName: "auth_request_dto.tmpl"},
			{Path: rolePath + "/interface/http/v1/dto/role_response.go", Package: "dto", TemplateName: "auth_response_dto.tmpl"},
			
			// Infrastructure layer
			{Path: rolePath + "/infrastructure/postgres/role_postgres.go", Package: "postgres", TemplateName: "role_postgres.tmpl"},
		}
		
		// Permission bounded context (permission management)
		permissionPath := "internal/permission"
		permissionFiles := []File{
			// Domain layer
			{Path: permissionPath + "/domain/entity/permission.go", Package: "entity", TemplateName: "permission_entity.tmpl"},
			{Path: permissionPath + "/domain/repository/permission_repository.go", Package: "repository", TemplateName: "permission_repository.tmpl"},
			
			// Application layer
			{Path: permissionPath + "/application/permission_service.go", Package: "application", TemplateName: "permission_service.tmpl"},
			
			// Interface layer
			{Path: permissionPath + "/interface/http/v1/handlers/permission_handler.go", Package: "handlers", TemplateName: "permission_handler.tmpl"},
			{Path: permissionPath + "/interface/http/v1/routes/permission_routes.go", Package: "routes", TemplateName: "permission_routes.tmpl"},
			{Path: permissionPath + "/interface/http/v1/dto/permission_request.go", Package: "dto", TemplateName: "auth_request_dto.tmpl"},
			{Path: permissionPath + "/interface/http/v1/dto/permission_response.go", Package: "dto", TemplateName: "permission_dto.tmpl"},
			
			// Infrastructure layer
			{Path: permissionPath + "/infrastructure/postgres/permission_postgres.go", Package: "postgres", TemplateName: "permission_postgres.tmpl"},
		}
		
		// Combine all auth-related files
		allAuthFiles := append(authUtilities, authFiles...)
		allAuthFiles = append(allAuthFiles, userFiles...)
		allAuthFiles = append(allAuthFiles, roleFiles...)
		allAuthFiles = append(allAuthFiles, permissionFiles...)
		
		boundedContextFiles = append(boundedContextFiles, allAuthFiles...)
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
		UseAuth:     fg.useAuth,
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
	
	// For auth-related templates in monolith mode, set specific import paths for separate bounded contexts
	if fg.isMonolith && fg.useAuth {
		// Auth service needs to import from user, role, and permission bounded contexts
		data.UserEntityImport = fmt.Sprintf("%s/internal/user/domain/entity", fg.moduleName)
		data.UserRepositoryImport = fmt.Sprintf("%s/internal/user/domain/repository", fg.moduleName)
		data.RoleEntityImport = fmt.Sprintf("%s/internal/role/domain/entity", fg.moduleName)
		data.RoleRepositoryImport = fmt.Sprintf("%s/internal/role/domain/repository", fg.moduleName)
		data.PermissionEntityImport = fmt.Sprintf("%s/internal/permission/domain/entity", fg.moduleName)
		data.PermissionRepositoryImport = fmt.Sprintf("%s/internal/permission/domain/repository", fg.moduleName)
		data.AuthServiceImport = fmt.Sprintf("%s/internal/auth/application", fg.moduleName)
		data.UserServiceImport = fmt.Sprintf("%s/internal/user/application", fg.moduleName)
		data.RoleServiceImport = fmt.Sprintf("%s/internal/role/application", fg.moduleName)
	}
	
	return data
}

// createFile creates a single file
func (fg *FileGenerator) createFile(file File, entityName string) error {
	fullPath := filepath.Join(fg.projectRoot, file.Path)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	// Prepare template data
	templateData := fg.prepareTemplateData(file.Package, entityName)
	
	// If a template is specified, render it
	if file.TemplateName != "" && entityName != "" {
		templatePath := "templates/" + file.TemplateName
		return fg.renderer.RenderToFile(templatePath, fullPath, templateData)
	} else if file.TemplateName == "db.tmpl" || file.TemplateName == "logger.tmpl" || 
		strings.HasPrefix(file.TemplateName, "auth_") ||
		strings.HasSuffix(file.TemplateName, "_entity.tmpl") ||
		strings.HasSuffix(file.TemplateName, "_service.tmpl") ||
		strings.HasSuffix(file.TemplateName, "_handler.tmpl") ||
		strings.HasSuffix(file.TemplateName, "_repository.tmpl") ||
		strings.HasSuffix(file.TemplateName, "_postgres.tmpl") ||
		file.TemplateName == "jwt_utils.tmpl" {
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
