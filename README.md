GoEasy
A CLI tool to generate Domain-Driven Design (DDD) Go service modules.
Usage
goeasy -generate <module-name>

Example:
goeasy -generate auth-service

Structure

/cmd/goeasy: CLI entry point
/internal/generator: Core generation logic
/internal/templates: File templates
