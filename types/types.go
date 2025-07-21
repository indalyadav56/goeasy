package types

type Config struct {
	ModuleName  string
	ProjectRoot string
	Entities    []string
}

type File struct {
	Path         string
	Package      string
	TemplateName string
}

type TemplateData struct {
	Package     string
	ProjectRoot string
	ModuleName  string
	EntityName  string
}
