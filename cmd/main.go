package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/indalyadav56/goeasy/internal/generator"
)

func main() {
	// Parse CLI arguments
	moduleName := flag.String("generate", "", "Module name (e.g., auth-service)")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("Error: Provide a module name using -generate")
		fmt.Println("Example: goeasy -generate auth-service")
		os.Exit(1)
	}

	// Generate module
	if err := generator.GenerateModule(*moduleName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated DDD module: %s\n", *moduleName)
}
