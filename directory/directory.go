package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirectoryManager struct{}

func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{}
}

func (dm *DirectoryManager) CreateProjectDirectories(projectRoot string) error {
	dirs := []string{
		"cmd/server",
		"internal/interface/http/v1/handlers",
		"internal/interface/http/v1/routes",
		"internal/interface/http/v1/dto",
		"internal/interface/grpc",
		"internal/application/services",
		"internal/domain/entity",
		"internal/domain/constants",
		"internal/domain/repository",
		"internal/infrastructure",
		"internal/infrastructure/postgres",
		"config",
		"pkg/logger",
		"pkg/db",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectRoot, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	return nil
}
