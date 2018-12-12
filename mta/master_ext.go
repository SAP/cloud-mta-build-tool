package mta

// EXT - mta extension schema
type EXT struct {
	// indicates MTA schema version, using semver.
	SchemaVersion *string `yaml:"_schema-version"`
	// A globally unique ID of this MTA extension. Unlimited string of unicode characters.
	ID string `yaml:"ID"`
	// A non-translatable description of this MTA extension. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	//  a globally unique ID of the MTA or the MTA extension which shall be extended by this descriptor
	Extends string `yaml:"extends"`
	// Application version, using semantic versioning standard
	Version string `yaml:"version,omitempty"`
	// The provider of this extension descriptor
	Provider string `yaml:"provider,omitempty"`
	// list of modules
	Modules []*ModuleExt `yaml:"modules,omitempty"`
	// Resource declarations. Resources can be anything required to run the application which is not provided by the application itself
	Resources []*ResourceExt `yaml:"resources,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

// ModuleExt - modules section in MTA extension
type ModuleExt struct {
	// An MTA internal module name. Names need to be unique within the MTA scope
	Name string
	// Provided property values can be accessed by "~{<name-of-provides-section>/<provided-property-name>}". Such expressions can be part of an arbitrary string
	Properties map[string]interface{} `yaml:"properties,omitempty"`
	// list of names either matching a resource name or a name provided by another module within the same MTA
	Requires []Requires `yaml:"requires,omitempty"`
	// List of provided names (MTA internal)to which properties (= configuration data) can be attached
	Provides []Provides `yaml:"provides,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor. Parameters are not made available to the module at runtime
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	// Build-parameters are specifically steering the behavior of build tools.
	BuildParams map[string]interface{} `yaml:"build-parameters,omitempty"`
}

// ResourceExt - can be anything required to run the application which is not provided by the application itself.
type ResourceExt struct {
	Name string
	// A type of a resource. This type is interpreted by and must be known to the deployer. Resources can be untyped
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	// property names and values make up the configuration data which is to be provided to requiring modules at runtime
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}
