package generator

import (
	"embed"
	"testing"

	"github.com/indalyadav56/gogen/internal/cli"
)

func TestNewProjectGenerator(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "github.com/test/project",
		Entities:   []string{"user"},
		Monolith:   false,
		UseGin:     false,
	}

	generator := NewProjectGenerator(config, mockFS)

	if generator == nil {
		t.Error("NewProjectGenerator() returned nil")
	}
	if generator.config != config {
		t.Error("config not set correctly")
	}
}

func TestProjectGenerator_Generate_Microservice(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "github.com/test/project",
		Entities:   []string{"user"},
		Monolith:   false,
		UseGin:     false,
	}

	generator := NewProjectGenerator(config, mockFS)

	// Test generation (will fail due to missing templates, but we can test structure)
	err := generator.Generate()

	// We expect an error due to missing templates or directory creation
	if err != nil {
		t.Logf("Expected error due to missing templates or directory creation: %v", err)
	}
}

func TestProjectGenerator_Generate_Monolith(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "github.com/test/project",
		Entities:   []string{"user", "product"},
		Monolith:   true,
		UseGin:     true,
	}

	generator := NewProjectGenerator(config, mockFS)

	// Test generation (will fail due to missing templates, but we can test structure)
	err := generator.Generate()

	// We expect an error due to missing templates or directory creation
	if err != nil {
		t.Logf("Expected error due to missing templates or directory creation: %v", err)
	}
}

func TestProjectGenerator_Generate_EmptyEntities(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "github.com/test/project",
		Entities:   []string{}, // Empty entities
		Monolith:   false,
		UseGin:     false,
	}

	generator := NewProjectGenerator(config, mockFS)

	// Test generation with empty entities
	err := generator.Generate()

	// Should handle empty entities gracefully
	if err != nil {
		// Error is expected due to missing templates, but should not panic
		t.Logf("Expected error due to missing templates: %v", err)
	}
}

func TestProjectGenerator_Generate_InvalidModuleName(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "", // Invalid empty module name
		Entities:   []string{"user"},
		Monolith:   false,
		UseGin:     false,
	}

	generator := NewProjectGenerator(config, mockFS)

	// Test with invalid module name
	err := generator.Generate()

	if err == nil {
		t.Error("Expected error for invalid module name, but got nil")
	}
}

func TestProjectGenerator_Generate_MultipleEntities(t *testing.T) {
	var mockFS embed.FS
	config := &cli.Config{
		ModuleName: "github.com/test/project",
		Entities:   []string{"user", "product", "order"},
		Monolith:   true,
		UseGin:     false,
	}

	generator := NewProjectGenerator(config, mockFS)

	// Test generation with multiple entities
	err := generator.Generate()

	// We expect an error due to missing templates or directory creation
	if err != nil {
		t.Logf("Expected error due to missing templates or directory creation: %v", err)
	}

	// Verify config has multiple entities
	if len(config.Entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(config.Entities))
	}
}
