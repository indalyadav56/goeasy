package cli

import (
	"flag"
	"strings"
)

// stringSlice implements flag.Value interface for handling multiple string flags
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Config holds all CLI configuration
type Config struct {
	ModuleName string
	Monolith   bool
	Entities   []string
	UseGin     bool
}

// ParseFlags parses command line flags and returns configuration
func ParseFlags() *Config {
	config := &Config{}
	
	moduleFlag := flag.String("module", "github.com/username/golang_project", "Go module name (e.g. github.com/user/project)")
	monolithFlag := flag.Bool("monolith", false, "for monolith architecture")
	ginFlag := flag.Bool("gin", false, "use Gin framework instead of Chi for HTTP routing")
	
	var entities stringSlice
	flag.Var(&entities, "entity", "Specify one or more entity names. Example: --entity User --entity Product")
	
	flag.Parse()
	
	config.ModuleName = *moduleFlag
	config.Monolith = *monolithFlag
	config.Entities = entities
	config.UseGin = *ginFlag
	
	return config
}

// GetProjectRoot extracts project root name from module name
func (c *Config) GetProjectRoot() string {
	moduleParts := strings.Split(c.ModuleName, "/")
	return moduleParts[len(moduleParts)-1]
}
