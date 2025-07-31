package main

import (
	"fmt"
	"os"

	goembed "github.com/indalyadav56/gogen"
	"github.com/indalyadav56/gogen/internal/cli"
	"github.com/indalyadav56/gogen/internal/generator"
)

func main() {
	// Parse command line flags
	config := cli.ParseFlags()
	
	// Create project generator with embedded templates
	projectGen := generator.NewProjectGenerator(config, goembed.TemplateFS)
	
	// Generate the project
	if err := projectGen.Generate(); err != nil {
		exitWithError(err)
	}
	
	fmt.Println("✅ Project structure scaffolded successfully.")
}

// exitWithError prints an error message and exits with status 1
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
	os.Exit(1)
}
