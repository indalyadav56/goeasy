package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirectoryStructure defines the project directory structure
type DirectoryStructure struct {
	ProjectRoot string
	EntityName  string
	IsMonolith  bool
}

// CreateDirectories creates the project directory structure
func (ds *DirectoryStructure) CreateDirectories() error {
	dirs := ds.getDirectories()
	
	for _, dir := range dirs {
		path := filepath.Join(ds.ProjectRoot, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	
	return nil
}

// getDirectories returns the list of directories to create based on architecture type
func (ds *DirectoryStructure) getDirectories() []string {
	if ds.IsMonolith {
		return ds.getMonolithDirectories()
	}
	return ds.getStandardDirectories()
}

// getStandardDirectories returns directories for standard architecture
func (ds *DirectoryStructure) getStandardDirectories() []string {
	return []string{
		"cmd",
		"docs",
		"docs/swagger",
		"docs/api",
		"internal/interface/http/v1/handlers",
		"internal/interface/http/v1/routes",
		"internal/interface/http/v1/dto",
		"internal/interface/http/middlewares",
		"internal/interface/grpc",
		"internal/application",
		"internal/domain/entity",
		"internal/domain/constants",
		"internal/domain/repository",
		"internal/infrastructure",
		"internal/infrastructure/postgres",
		"migrations",
		"config",
		"pkg/logger",
		"pkg/db",
	}
}

// getMonolithDirectories returns directories for monolith architecture with bounded contexts
func (ds *DirectoryStructure) getMonolithDirectories() []string {
	// Base directories
	dirs := []string{
		"cmd",
		"docs",
		"docs/swagger",
		"docs/api",
		"config",
		"pkg/logger",
		"pkg/db",
		"internal/shared/middleware",
		"internal/shared/dto",
		"internal/shared/utils",
		"migrations",
		"scripts",
	}
	
	// If no specific entity, create a generic structure
	if ds.EntityName == "" {
		return append(dirs, []string{
			"internal/example/domain",
			"internal/example/repository",
			"internal/example/service",
			"internal/example/handler",
			"internal/example/infrastructure",
		}...)
	}
	
	// Create bounded context structure for the entity using clean architecture
	entityPath := fmt.Sprintf("internal/%s", strings.ToLower(ds.EntityName))
	boundedContextDirs := []string{
		// Domain layer - core business logic (entities, value objects, aggregates)
		entityPath + "/domain",
		entityPath + "/domain/constants",
		entityPath + "/domain/entity",
		entityPath + "/domain/repository",
		
		// Interface layer - interface adapters (HTTP handlers, controllers)
		entityPath + "/interface",
		entityPath + "/interface/http",
		entityPath + "/interface/grpc",
		entityPath + "/interface/http/v1",
		entityPath + "/interface/http/v1/handlers",
		entityPath + "/interface/http/v1/routes",
		entityPath + "/interface/http/v1/dto",
		entityPath + "/interface/http/middlewares",
		
		// Application layer - application and use cases
		entityPath + "/application",
		
		// Infrastructure layer - external concerns (database, external APIs)
		entityPath + "/infrastructure",
		entityPath + "/infrastructure/postgres",
	}
	
	return append(dirs, boundedContextDirs...)
}
