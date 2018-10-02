package mta

// BuildParameters - build params
type BuildParameters struct {
	Builder  string          `yaml:"builder,omitempty"`
	Type     string          `yaml:"type,omitempty"`
	Path     string          `yaml:"path,omitempty"`
	Requires []BuildRequires `yaml:"requires,omitempty"`
}
