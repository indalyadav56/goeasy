package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryStructure_CreateDirectories_Microservice(t *testing.T) {
	tempDir := t.TempDir()
	ds := &DirectoryStructure{
		ProjectRoot: tempDir,
		EntityName:  "user",
		IsMonolith:  false,
	}

	err := ds.CreateDirectories()
	if err != nil {
		t.Fatalf("CreateDirectories() error = %v", err)
	}

	// Verify some key directories were created
	expectedDirs := []string{
		"cmd",
		"config",
		"pkg/logger",
		"internal/application",
		"internal/domain/entity",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Directory %s was not created", fullPath)
		}
	}
}

func TestDirectoryStructure_GetDirectories_Microservice(t *testing.T) {
	ds := &DirectoryStructure{
		ProjectRoot: "/test/project",
		EntityName:  "user",
		IsMonolith:  false,
	}
	dirs := ds.getDirectories()

	// Check that microservice directories are included
	expectedDirs := []string{
		"cmd",
		"config",
		"pkg/logger",
		"pkg/db",
		"internal/interface/http/v1/handlers",
		"internal/interface/http/v1/routes",
		"internal/interface/http/v1/dto",
		"internal/interface/http/middlewares",
		"internal/interface/grpc",
		"internal/application",
		"internal/domain/entity",
		"internal/domain/constants",
		"internal/domain/repository",
		"internal/infrastructure/postgres",
		"migrations",
	}

	for _, expectedDir := range expectedDirs {
		if !contains(dirs, expectedDir) {
			t.Errorf("Expected directory %s not found in microservice structure", expectedDir)
		}
	}
}

func TestDirectoryStructure_GetDirectories_Monolith(t *testing.T) {
	ds := &DirectoryStructure{
		ProjectRoot: "/test/project",
		EntityName:  "user",
		IsMonolith:  true,
	}
	dirs := ds.getDirectories()

	// Check that monolith base directories are included
	expectedBaseDirs := []string{
		"cmd",
		"config",
		"pkg/logger",
		"pkg/db",
		"internal/shared/middleware",
		"internal/shared/dto",
		"internal/shared/utils",
		"migrations",
		"scripts",
	}

	// Check that user bounded context directories are included
	expectedUserDirs := []string{
		"internal/user/domain",
		"internal/user/domain/constants",
		"internal/user/domain/entity",
		"internal/user/domain/repository",
		"internal/user/interface",
		"internal/user/interface/grpc",
		"internal/user/interface/http/v1",
		"internal/user/interface/http/v1/handlers",
		"internal/user/interface/http/v1/routes",
		"internal/user/interface/http/v1/dto",
		"internal/user/interface/http/middlewares",
		"internal/user/application",
		"internal/user/infrastructure",
		"internal/user/infrastructure/postgres",
	}

	allExpectedDirs := append(expectedBaseDirs, expectedUserDirs...)

	for _, expectedDir := range allExpectedDirs {
		if !contains(dirs, expectedDir) {
			t.Errorf("Expected directory %s not found in monolith structure", expectedDir)
		}
	}
}

func TestDirectoryStructure_GetDirectories_MonolithMultipleEntities(t *testing.T) {
	// Test with multiple entities by creating separate DirectoryStructure instances
	userDS := &DirectoryStructure{
		ProjectRoot: "/test/project",
		EntityName:  "user",
		IsMonolith:  true,
	}
	productDS := &DirectoryStructure{
		ProjectRoot: "/test/project",
		EntityName:  "product",
		IsMonolith:  true,
	}

	userDirs := userDS.getDirectories()
	productDirs := productDS.getDirectories()

	// Check user-specific directories
	expectedUserDirs := []string{
		"internal/user/domain",
		"internal/user/application",
		"internal/user/interface/http/v1/handlers",
	}

	// Check product-specific directories
	expectedProductDirs := []string{
		"internal/product/domain",
		"internal/product/application",
		"internal/product/interface/http/v1/handlers",
	}

	for _, expectedDir := range expectedUserDirs {
		if !contains(userDirs, expectedDir) {
			t.Errorf("Expected user directory %s not found", expectedDir)
		}
	}

	for _, expectedDir := range expectedProductDirs {
		if !contains(productDirs, expectedDir) {
			t.Errorf("Expected product directory %s not found", expectedDir)
		}
	}
}

func TestDirectoryStructure_CreateDirectories_Monolith(t *testing.T) {
	tempDir := t.TempDir()
	ds := &DirectoryStructure{
		ProjectRoot: tempDir,
		EntityName:  "user",
		IsMonolith:  true,
	}

	err := ds.CreateDirectories()
	if err != nil {
		t.Fatalf("CreateDirectories() error = %v", err)
	}

	// Verify bounded context directories were created
	expectedDirs := []string{
		"internal/user/domain",
		"internal/user/application",
		"internal/user/interface/http/v1/handlers",
		"internal/user/infrastructure/postgres",
		"internal/shared/middleware",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Bounded context directory %s was not created", fullPath)
		}
	}
}

func TestDirectoryStructure_GetStandardDirectories(t *testing.T) {
	ds := &DirectoryStructure{
		ProjectRoot: "/test",
		EntityName:  "user",
		IsMonolith:  false,
	}

	dirs := ds.getStandardDirectories()

	// Check that all expected microservice directories are present
	expectedDirs := []string{
		"cmd",
		"config",
		"pkg/logger",
		"pkg/db",
		"internal/application",
		"internal/domain/entity",
		"internal/infrastructure/postgres",
		"migrations",
	}

	for _, expected := range expectedDirs {
		if !contains(dirs, expected) {
			t.Errorf("Expected microservice directory %s not found", expected)
		}
	}
}

func TestDirectoryStructure_GetMonolithDirectories(t *testing.T) {
	ds := &DirectoryStructure{
		ProjectRoot: "/test",
		EntityName:  "user",
		IsMonolith:  true,
	}

	dirs := ds.getMonolithDirectories()

	// Check base directories
	expectedBaseDirs := []string{
		"cmd",
		"config",
		"pkg/logger",
		"pkg/db",
		"internal/shared/middleware",
		"migrations",
		"scripts",
	}

	// Check user bounded context directories
	expectedUserDirs := []string{
		"internal/user/domain",
		"internal/user/application",
		"internal/user/interface",
		"internal/user/infrastructure",
	}

	allExpected := append(expectedBaseDirs, expectedUserDirs...)

	for _, expected := range allExpected {
		if !contains(dirs, expected) {
			t.Errorf("Expected monolith directory %s not found", expected)
		}
	}
}

func TestDirectoryStructure_GetMonolithDirectories_EmptyEntity(t *testing.T) {
	ds := &DirectoryStructure{
		ProjectRoot: "/test",
		EntityName:  "",
		IsMonolith:  true,
	}

	dirs := ds.getMonolithDirectories()

	// Should create example bounded context when no entity specified
	expectedExampleDirs := []string{
		"internal/example/domain",
		"internal/example/repository",
		"internal/example/service",
		"internal/example/handler",
		"internal/example/infrastructure",
	}

	for _, expected := range expectedExampleDirs {
		if !contains(dirs, expected) {
			t.Errorf("Expected example directory %s not found when entity is empty", expected)
		}
	}
}



// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}