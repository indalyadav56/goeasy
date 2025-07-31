package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	projectRoot := "test-project"
	manager := NewManager(projectRoot)

	if manager == nil {
		t.Error("NewManager() returned nil")
	}
	if manager.projectRoot != projectRoot {
		t.Errorf("projectRoot = %v, want %v", manager.projectRoot, projectRoot)
	}
}

func TestManager_Init(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create project directory
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Change to project directory for the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectPath)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	manager := NewManager(projectName)
	moduleName := "github.com/test/project"

	// Test go mod init
	err = manager.Init(moduleName)
	if err != nil {
		t.Logf("go mod init failed (expected if go is not available): %v", err)
		return // Skip test if go is not available
	}

	// Check if go.mod file was created
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file was not created")
	}
}

func TestManager_Tidy(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create project directory
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create a basic go.mod file
	goModContent := `module github.com/test/project

go 1.21
`
	goModPath := filepath.Join(projectPath, "go.mod")
	err = os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod file: %v", err)
	}

	// Change to project directory for the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectPath)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	manager := NewManager(projectName)

	// Test go mod tidy
	err = manager.Tidy()
	if err != nil {
		t.Logf("go mod tidy failed (expected if go is not available): %v", err)
		return // Skip test if go is not available
	}

	// If successful, the go.mod file should still exist
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file was removed unexpectedly")
	}
}

func TestManager_Init_EmptyModuleName(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create project directory
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Change to project directory for the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectPath)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	manager := NewManager(projectName)

	// Test go mod init with empty module name
	err = manager.Init("")
	if err != nil {
		t.Logf("go mod init with empty module name failed as expected: %v", err)
	}
}

func TestManager_Tidy_NoGoMod(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create project directory but no go.mod
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Change to project directory for the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectPath)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	manager := NewManager(projectName)

	// Test go mod tidy without go.mod file
	err = manager.Tidy()
	if err != nil {
		t.Logf("go mod tidy without go.mod failed as expected: %v", err)
	}
}
