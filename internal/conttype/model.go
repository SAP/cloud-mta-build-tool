package conttype

// ContentTypes - data structure for content types configuration
type ContentTypes struct {
	ContentTypes []ContentType `yaml:"content-types"`
}

// ContentType -  content types
type ContentType struct {
	Extension   string `yaml:"extension"`
	ContentType string `yaml:"content-type"`
}
