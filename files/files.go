package files

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// type FileGenerator struct {
// 	renderer *TemplateRenderer
// }

// func NewFileGenerator(renderer *TemplateRenderer) *FileGenerator {
// 	return &FileGenerator{
// 		renderer: renderer,
// 	}
// }

// func (fg *FileGenerator) GetFileDefinitions() []File {
// 	return []File{
// 		{Path: "cmd/server/main.go", Package: "main", TemplateName: "main.tmpl"},
// 		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
// 		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},

// 		// Domain
// 		{Path: "internal/domain/constants/constants.go", Package: "constants", TemplateName: ""},
// 		{Path: "internal/domain/entity/entity.go", Package: "entity", TemplateName: "entity.tmpl"},
// 		{Path: "internal/domain/repository/repository.go", Package: "repository", TemplateName: "repository.tmpl"},

// 		{Path: "internal/application/services/service.go", Package: "services", TemplateName: "service.tmpl"},
// 		{Path: "internal/interface/http/v1/routes/routes.go", Package: "routes", TemplateName: "routes.tmpl"},
// 		{Path: "internal/interface/http/v1/handlers/handler.go", Package: "handlers", TemplateName: "handler.tmpl"},
// 		{Path: "internal/infrastructure/postgres/postgres.go", Package: "postgres", TemplateName: "postgres_repository.tmpl"},

// 		// DTO
// 		{Path: "internal/interface/http/v1/dto/request.go", Package: "dto", TemplateName: ""},
// 		{Path: "internal/interface/http/v1/dto/response.go", Package: "dto", TemplateName: ""},

// 		{Path: "pkg/db/db.go", Package: "db", TemplateName: "db.tmpl"},
// 	}
// }

// func (fg *FileGenerator) GenerateFiles(config *Config, entityName string) error {
// 	files := fg.GetFileDefinitions()

// 	for _, f := range files {
// 		fullPath := filepath.Join(config.ProjectRoot, f.Path)

// 		// Prepare template data
// 		templateData := TemplateData{
// 			Package:     f.Package,
// 			ProjectRoot: config.ProjectRoot,
// 			ModuleName:  config.ModuleName,
// 			EntityName:  entityName,
// 		}

// 		// If a template is specified, render it
// 		if f.TemplateName != "" && entityName != "" {
// 			templatePath := filepath.Join("templates", f.TemplateName)
// 			if err := fg.renderer.RenderToFile(templatePath, fullPath, templateData); err != nil {
// 				return fmt.Errorf("failed to render template %s: %w", f.TemplateName, err)
// 			}
// 		} else if f.TemplateName == "db.tmpl" || f.TemplateName == "logger.tmpl" {
// 			templatePath := filepath.Join("templates", f.TemplateName)
// 			if err := fg.renderer.RenderToFile(templatePath, fullPath, templateData); err != nil {
// 				return fmt.Errorf("failed to render template %s: %w", f.TemplateName, err)
// 			}
// 		} else {
// 			content := fmt.Sprintf("package %s\n", f.Package)
// 			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
// 				return fmt.Errorf("failed to create file %s: %w", fullPath, err)
// 			}
// 		}
// 	}

// 	return nil
// }
