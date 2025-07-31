package cli

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "default values",
			args: []string{},
			expected: Config{
				ModuleName: "github.com/username/golang_project",
				Monolith:   false,
				Entities:   []string{},
				UseGin:     false,
			},
		},
		{
			name: "custom module name",
			args: []string{"--module", "github.com/test/project"},
			expected: Config{
				ModuleName: "github.com/test/project",
				Monolith:   false,
				Entities:   []string{},
				UseGin:     false,
			},
		},
		{
			name: "monolith flag",
			args: []string{"--monolith"},
			expected: Config{
				ModuleName: "github.com/username/golang_project",
				Monolith:   true,
				Entities:   []string{},
				UseGin:     false,
			},
		},
		{
			name: "gin flag",
			args: []string{"--gin"},
			expected: Config{
				ModuleName: "github.com/username/golang_project",
				Monolith:   false,
				Entities:   []string{},
				UseGin:     true,
			},
		},
		{
			name: "single entity",
			args: []string{"--entity", "user"},
			expected: Config{
				ModuleName: "github.com/username/golang_project",
				Monolith:   false,
				Entities:   []string{"user"},
				UseGin:     false,
			},
		},
		{
			name: "multiple entities",
			args: []string{"--entity", "user", "--entity", "product"},
			expected: Config{
				ModuleName: "github.com/username/golang_project",
				Monolith:   false,
				Entities:   []string{"user", "product"},
				UseGin:     false,
			},
		},
		{
			name: "all flags combined",
			args: []string{
				"--module", "github.com/company/api",
				"--entity", "user",
				"--entity", "order",
				"--monolith",
				"--gin",
			},
			expected: Config{
				ModuleName: "github.com/company/api",
				Monolith:   true,
				Entities:   []string{"user", "order"},
				UseGin:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			
			// Set os.Args to simulate command line input
			oldArgs := os.Args
			os.Args = append([]string{"gogen"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			config := ParseFlags()

			if !reflect.DeepEqual(config, tt.expected) {
				t.Errorf("ParseFlags() = %+v, want %+v", config, tt.expected)
			}
		})
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "empty slice",
			values:   []string{},
			expected: "",
		},
		{
			name:     "single value",
			values:   []string{"user"},
			expected: "user",
		},
		{
			name:     "multiple values",
			values:   []string{"user", "product", "order"},
			expected: "user,product,order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ss stringSlice
			for _, v := range tt.values {
				ss.Set(v)
			}

			result := ss.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}

			// Test that the slice contains expected values
			if !reflect.DeepEqual([]string(ss), tt.values) {
				t.Errorf("stringSlice = %v, want %v", []string(ss), tt.values)
			}
		})
	}
}
