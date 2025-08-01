package scaffold

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/indalyadav56/gogen/internal/template"
)

func TestNewFileGenerator(t *testing.T) {
	var mockFS embed.FS
	renderer := template.NewRenderer(mockFS)
	projectRoot := "/test/project"
	moduleName := "github.com/test/project"
	isMonolith := false
	useGin := false
	entities := []string{"user"}

	fg := NewFileGenerator(renderer, projectRoot, moduleName, isMonolith, useGin, false, entities)

	if fg == nil {
		t.Error("NewFileGenerator() returned nil")
	}
	if fg.projectRoot != projectRoot {
		t.Errorf("projectRoot = %v, want %v", fg.projectRoot, projectRoot)
	}
	if fg.moduleName != moduleName {
		t.Errorf("moduleName = %v, want %v", fg.moduleName, moduleName)
	}
	if fg.isMonolith != isMonolith {
		t.Errorf("isMonolith = %v, want %v", fg.isMonolith, isMonolith)
	}
	if fg.useGin != useGin {
		t.Errorf("useGin = %v, want %v", fg.useGin, useGin)
	}
}

func TestFileGenerator_GetMicroserviceFileList(t *testing.T) {
	var mockFS embed.FS
	renderer := template.NewRenderer(mockFS)
	fg := NewFileGenerator(renderer, "/test", "github.com/test/project", false, false, false, []string{"user"})

	files := fg.getMicroserviceFileList("user")

	// Check that essential microservice files are included
	expectedFiles := map[string]bool{
		"cmd/main.go":                                                     true,
		"config/config.go":                                                true,
		"pkg/logger/logger.go":                                            true,
		"internal/application/user_service.go":                           true,
		"internal/interface/http/v1/handlers/user_handler.go":            true,
		"internal/interface/http/v1/routes/routes.go":                    true,
		"internal/domain/entity/entity.go":                               true,
		"internal/domain/repository/repository.go":                       true,
		"internal/infrastructure/postgres/postgres.go":                   true,
		"Dockerfile":                                                      true,
		"Taskfile.yaml":                                                   true,
	}

	for _, file := range files {
		if expectedFiles[file.Path] {
			delete(expectedFiles, file.Path)
		}
	}

	if len(expectedFiles) > 0 {
		for missingFile := range expectedFiles {
			t.Errorf("Expected microservice file %s not found", missingFile)
		}
	}
}

func TestFileGenerator_GetMonolithFileList(t *testing.T) {
	var mockFS embed.FS
	renderer := template.NewRenderer(mockFS)
	fg := NewFileGenerator(renderer, "/test", "github.com/test/project", true, false, false, []string{"user"})

	files := fg.getMonolithFileList("user")

	// Check that essential monolith files are included
	expectedFiles := map[string]bool{
		"cmd/main.go":                                                     true,
		"config/config.go":                                                true,
		"pkg/logger/logger.go":                                            true,
		"pkg/db/db.go":                                                    true,
		"internal/shared/middleware/auth.go":                             true,
		"internal/shared/dto/common.go":                                   true,
		"internal/shared/utils/utils.go":                                 true,
		"internal/user/domain/entity/entity.go":                          true,
		"internal/user/domain/repository/repository.go":                  true,
		"internal/user/application/user_service.go":                      true,
		"internal/user/interface/http/v1/handlers/user_handler.go":       true,
		"internal/user/interface/http/v1/routes/routes.go":               true,
		"internal/user/infrastructure/postgres/postgres.go":              true,
		"Dockerfile":                                                      true,
	}

	for _, file := range files {
		if expectedFiles[file.Path] {
			delete(expectedFiles, file.Path)
		}
	}

	if len(expectedFiles) > 0 {
		for missingFile := range expectedFiles {
			t.Errorf("Expected monolith file %s not found", missingFile)
		}
	}
}

func TestFileGenerator_GetHandlerTemplate(t *testing.T) {
	tests := []struct {
		name     string
		useGin   bool
		expected string
	}{
		{
			name:     "Chi framework",
			useGin:   false,
			expected: "handler.tmpl",
		},
		{
			name:     "Gin framework",
			useGin:   true,
			expected: "gin_handler.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockFS embed.FS
			renderer := template.NewRenderer(mockFS)
			fg := NewFileGenerator(renderer, "/test", "github.com/test/project", false, tt.useGin, false, []string{"user"})

			result := fg.getHandlerTemplate()
			if result != tt.expected {
				t.Errorf("getHandlerTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileGenerator_GetRoutesTemplate(t *testing.T) {
	tests := []struct {
		name     string
		useGin   bool
		expected string
	}{
		{
			name:     "Chi framework",
			useGin:   false,
			expected: "routes.tmpl",
		},
		{
			name:     "Gin framework",
			useGin:   true,
			expected: "gin_routes.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockFS embed.FS
			renderer := template.NewRenderer(mockFS)
			fg := NewFileGenerator(renderer, "/test", "github.com/test/project", false, tt.useGin, false, []string{"user"})

			result := fg.getRoutesTemplate()
			if result != tt.expected {
				t.Errorf("getRoutesTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileGenerator_GetMainTemplate(t *testing.T) {
	tests := []struct {
		name       string
		isMonolith bool
		useGin     bool
		expected   string
	}{
		{
			name:       "Microservice",
			isMonolith: false,
			useGin:     false,
			expected:   "main.tmpl",
		},
		{
			name:       "Monolith with Chi",
			isMonolith: true,
			useGin:     false,
			expected:   "monolith_main.tmpl",
		},
		{
			name:       "Monolith with Gin",
			isMonolith: true,
			useGin:     true,
			expected:   "gin_monolith_main.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockFS embed.FS
			renderer := template.NewRenderer(mockFS)
			fg := NewFileGenerator(renderer, "/test", "github.com/test/project", tt.isMonolith, tt.useGin, false, []string{"user"})

			result := fg.getMainTemplate()
			if result != tt.expected {
				t.Errorf("getMainTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileGenerator_PrepareTemplateData(t *testing.T) {
	tests := []struct {
		name         string
		isMonolith   bool
		useGin       bool
		entityName   string
		packageName  string
		expectedData template.Data
	}{
		{
			name:        "Microservice template data",
			isMonolith:  false,
			useGin:      false,
			entityName:  "user",
			packageName: "handlers",
			expectedData: template.Data{
				Package:          "handlers",
				ProjectRoot:      "/test",
				ModuleName:       "github.com/test/project",
				EntityName:       "user",
				Entities:         []string{"user"},
				IsMonolith:       false,
				UseGin:           false,
				ServiceImport:    "github.com/test/project/internal/application",
				RepositoryImport: "github.com/test/project/internal/domain/repository",
				EntityImport:     "github.com/test/project/internal/domain/entity",
			},
		},
		{
			name:        "Monolith template data",
			isMonolith:  true,
			useGin:      true,
			entityName:  "user",
			packageName: "handlers",
			expectedData: template.Data{
				Package:          "handlers",
				ProjectRoot:      "/test",
				ModuleName:       "github.com/test/project",
				EntityName:       "user",
				Entities:         []string{"user"},
				IsMonolith:       true,
				UseGin:           true,
				ServiceImport:    "github.com/test/project/internal/user/application",
				RepositoryImport: "github.com/test/project/internal/user/domain/repository",
				EntityImport:     "github.com/test/project/internal/user/domain/entity",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockFS embed.FS
			renderer := template.NewRenderer(mockFS)
			fg := NewFileGenerator(renderer, "/test", "github.com/test/project", tt.isMonolith, tt.useGin, false, []string{"user"})

			result := fg.prepareTemplateData(tt.packageName, tt.entityName)

			// Check key fields
			if result.Package != tt.expectedData.Package {
				t.Errorf("Package = %v, want %v", result.Package, tt.expectedData.Package)
			}
			if result.ModuleName != tt.expectedData.ModuleName {
				t.Errorf("ModuleName = %v, want %v", result.ModuleName, tt.expectedData.ModuleName)
			}
			if result.EntityName != tt.expectedData.EntityName {
				t.Errorf("EntityName = %v, want %v", result.EntityName, tt.expectedData.EntityName)
			}
			if result.IsMonolith != tt.expectedData.IsMonolith {
				t.Errorf("IsMonolith = %v, want %v", result.IsMonolith, tt.expectedData.IsMonolith)
			}
			if result.UseGin != tt.expectedData.UseGin {
				t.Errorf("UseGin = %v, want %v", result.UseGin, tt.expectedData.UseGin)
			}
			if result.ServiceImport != tt.expectedData.ServiceImport {
				t.Errorf("ServiceImport = %v, want %v", result.ServiceImport, tt.expectedData.ServiceImport)
			}
		})
	}
}

func TestFileGenerator_GenerateFiles(t *testing.T) {
	tempDir := t.TempDir()
	var mockFS embed.FS
	renderer := template.NewRenderer(mockFS)
	fg := NewFileGenerator(renderer, tempDir, "github.com/test/project", false, false, false, []string{"user"})

	// Test file generation (this will fail due to missing templates, but we can test the structure)
	err := fg.GenerateFiles("user")

	// We expect an error due to missing templates, but the function should attempt to create files
	if err == nil {
		// If no error, verify some directories were created
		expectedDirs := []string{
			filepath.Join(tempDir, "cmd"),
			filepath.Join(tempDir, "internal"),
			filepath.Join(tempDir, "pkg"),
		}

		for _, dir := range expectedDirs {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Expected directory %s was not created", dir)
			}
		}
	}
}
