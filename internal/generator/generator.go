package generator

import (
	"embed"
	"fmt"

	"github.com/indalyadav56/gogen/internal/cli"
	"github.com/indalyadav56/gogen/internal/gomod"
	"github.com/indalyadav56/gogen/internal/scaffold"
	"github.com/indalyadav56/gogen/internal/template"
	"github.com/indalyadav56/gogen/utils"
)

// ProjectGenerator handles the entire project generation process
type ProjectGenerator struct {
	config     *cli.Config
	renderer   *template.Renderer
	gomodMgr   *gomod.Manager
	projectRoot string
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator(config *cli.Config, templateFS embed.FS) *ProjectGenerator {
	projectRoot := config.GetProjectRoot()
	
	return &ProjectGenerator{
		config:      config,
		renderer:    template.NewRenderer(templateFS),
		gomodMgr:    gomod.NewManager(projectRoot),
		projectRoot: projectRoot,
	}
}

// Generate generates the complete project structure
func (pg *ProjectGenerator) Generate() error {
	// For monolith architecture, create directories for each entity
	if pg.config.Monolith && len(pg.config.Entities) > 0 {
		for _, entityName := range pg.config.Entities {
			if err := pg.createDirectoriesForEntity(utils.ToCamelCase(entityName)); err != nil {
				return fmt.Errorf("failed to create directories for entity %s: %w", entityName, err)
			}
		}
	} else {
		// Create base directory structure
		if err := pg.createDirectories(); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
	}
	
	// Generate files for entities
	if err := pg.generateFiles(); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}
	
	// Initialize Go module
	if err := pg.gomodMgr.Init(pg.config.ModuleName); err != nil {
		return err
	}
	
	// Run go mod tidy
	if err := pg.gomodMgr.Tidy(); err != nil {
		return err
	}
	
	return nil
}

// createDirectories creates the project directory structure
func (pg *ProjectGenerator) createDirectories() error {
	dirStructure := &scaffold.DirectoryStructure{
		ProjectRoot: pg.projectRoot,
		EntityName:  "",
		IsMonolith:  pg.config.Monolith,
	}
	
	return dirStructure.CreateDirectories()
}

// createDirectoriesForEntity creates directories for a specific entity in monolith architecture
func (pg *ProjectGenerator) createDirectoriesForEntity(entityName string) error {
	dirStructure := &scaffold.DirectoryStructure{
		ProjectRoot: pg.projectRoot,
		EntityName:  entityName,
		IsMonolith:  pg.config.Monolith,
	}
	
	return dirStructure.CreateDirectories()
}

// generateFiles generates all project files
func (pg *ProjectGenerator) generateFiles() error {
	fileGenerator := scaffold.NewFileGenerator(pg.renderer, pg.projectRoot, pg.config.ModuleName, pg.config.Monolith, pg.config.UseGin, pg.config.UseAuth, pg.config.Entities)
	
	if len(pg.config.Entities) > 0 {
		for _, entityName := range pg.config.Entities {
			if err := fileGenerator.GenerateFiles(utils.ToCamelCase(entityName)); err != nil {
				return err
			}
		}
	} else {
		// Generate files without specific entity
		if err := fileGenerator.GenerateFiles(""); err != nil {
			return err
		}
	}
	
	return nil
}
