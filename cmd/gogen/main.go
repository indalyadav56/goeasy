package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	goembed "github.com/indalyadav56/gogen"
	"github.com/indalyadav56/gogen/utils"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type File struct {
	Path         string
	Package      string
	TemplateName string
}

type TemplateData struct {
	Package     string
	ProjectRoot string
	ModuleName  string
	EntityName  string
}

var templateFS embed.FS

func main() {
	templateFS = goembed.TemplateFS

	moduleFlag := flag.String("module", "github.com/username/golang_project", "Go module name (e.g. github.com/user/project)")

	var entities stringSlice

	flag.Var(&entities, "entity", "Specify one or more entity names. Example: --entity User --entity Product")

	flag.Parse()

	var entityName string

	moduleName := *moduleFlag

	moduleParts := strings.Split(moduleName, "/")
	projectRoot := moduleParts[len(moduleParts)-1]

	// Define directories
	dirs := []string{
		"cmd/server",
		"docs",
		"docs/swagger",
		"docs/api",
		"internal/interface/http/v1/handlers",
		"internal/interface/http/v1/routes",
		"internal/interface/http/v1/dto",
		"internal/interface/http/middlewares",
		"internal/interface/grpc",
		"internal/application/services",
		"internal/domain/entity",
		"internal/domain/constants",
		"internal/domain/repository",
		"internal/infrastructure",
		"internal/infrastructure/postgres",
		"internal/infrastructure/postgres/migrations",
		"config",
		"pkg/logger",
		"pkg/db",
	}

	// Create folders
	for _, dir := range dirs {
		path := filepath.Join(projectRoot, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			exitWithError(fmt.Errorf("failed to create directory: %w", err))
		}
	}

	if len(entities) > 0 {
		for _, entityName := range entities {
			createStructure(utils.ToCamelCase(entityName), moduleName, projectRoot)
		}
	} else {
		createStructure(entityName, moduleName, projectRoot)
	}

	// Init Go module
	if err := runGoModInit(moduleName, projectRoot); err != nil {
		exitWithError(fmt.Errorf("go mod init failed: %w", err))
	}

	if err := runGoModTidy(projectRoot); err != nil {
		exitWithError(fmt.Errorf("go mod tidy failed: %w", err))
	}

	fmt.Println("✅ Project structure scaffolded successfully.")
}

func runGoModInit(moduleName, projectRoot string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
	os.Exit(1)
}

func renderTemplateToFile(templatePath, outputPath string, data any) error {
	// Define custom template functions
	funcMap := template.FuncMap{
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"ToCamelCase": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		"ToPascalCase": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
	}

	// Read template content from embedded filesystem
	templateContent, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded template file %s: %w", templatePath, err)
	}

	// Parse template from content with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template content: %w", err)
	}

	// Create file to write rendered content
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func runGoModTidy(projectRoot string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createStructure(entityName, moduleName, projectRoot string) {
	files := []File{
		{Path: "cmd/server/main.go", Package: "main", TemplateName: "main.tmpl"},
		{Path: "config/config.go", Package: "config", TemplateName: "config.tmpl"},
		{Path: "pkg/logger/logger.go", Package: "logger", TemplateName: "logger.tmpl"},

		// Domain
		{Path: "internal/domain/constants/constants.go", Package: "constants", TemplateName: ""},
		{Path: "internal/domain/entity/entity.go", Package: "entity", TemplateName: "entity.tmpl"},
		{Path: "internal/domain/repository/repository.go", Package: "repository", TemplateName: "repository.tmpl"},

		{Path: "internal/application/services/" + fmt.Sprintf("%s_service.go", strings.ToLower(entityName)), Package: "services", TemplateName: "service.tmpl"},
		{Path: "internal/interface/http/v1/routes/routes.go", Package: "routes", TemplateName: "routes.tmpl"},
		{Path: "internal/interface/http/v1/handlers/" + fmt.Sprintf("%s_handler.go", strings.ToLower(entityName)), Package: "handlers", TemplateName: "handler.tmpl"},
		{Path: "internal/interface/http/middlewares/auth_middleware.go", Package: "middlewares"},

		{Path: "internal/infrastructure/postgres/postgres.go", Package: "postgres", TemplateName: "postgres_repository.tmpl"},

		// DTO
		{Path: "internal/interface/http/v1/dto/request.go", Package: "dto", TemplateName: ""},
		{Path: "internal/interface/http/v1/dto/response.go", Package: "dto", TemplateName: ""},

		{Path: "pkg/db/db.go", Package: "db", TemplateName: "db.tmpl"},

		{Path: ".gitignore", Package: "", TemplateName: ""},
		{Path: "Dockerfile", Package: "", TemplateName: ""},
	}

	// Create files
	for _, f := range files {
		fullPath := filepath.Join(projectRoot, f.Path)

		// Prepare template data
		templateData := TemplateData{
			Package:     f.Package,
			ProjectRoot: projectRoot,
			ModuleName:  moduleName,
			EntityName:  entityName,
		}

		// If a template is specified, render it
		if f.TemplateName != "" && entityName != "" {
			templatePath := filepath.Join("templates", f.TemplateName)
			if err := renderTemplateToFile(templatePath, fullPath, templateData); err != nil {
				exitWithError(fmt.Errorf("failed to render template %s: %w", f.TemplateName, err))
			}
		} else if f.TemplateName == "db.tmpl" || f.TemplateName == "logger.tmpl" {
			templatePath := filepath.Join("templates", f.TemplateName)
			if err := renderTemplateToFile(templatePath, fullPath, templateData); err != nil {
				exitWithError(fmt.Errorf("failed to render template %s: %w", f.TemplateName, err))
			}
		} else {
			content := fmt.Sprintf("package %s\n", f.Package)
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				exitWithError(fmt.Errorf("failed to create file %s: %w", fullPath, err))
			}
		}
	}
}
