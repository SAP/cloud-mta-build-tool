package mta

// MTA master schema, the schema will contain the latest mta schema version
// and all the previous version will be as subset of the latest
// Todo - Add the missing properties to support the latest 3.2 version
type MTA struct {
	// indicates MTA schema version, using semver.
	SchemaVersion *string `yaml:"_schema-version"`
	// A globally unique ID of this MTA. Unlimited string of unicode characters.
	ID string `yaml:"ID"`
	// A non-translatable description of this MTA. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	// Application version, using semantic versioning standard
	Version string `yaml:"version,omitempty"`
	// The provider or vendor of this software
	Provider string `yaml:"provider,omitempty"`
	// A copyright statement from the provider
	Copyright string `yaml:"copyright,omitempty"`
	// list of modules
	Modules []*Module `yaml:"modules,omitempty"`
	// Resource declarations. Resources can be anything required to run the application which is not provided by the application itself
	Resources []*Resource `yaml:"resources,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

// Module - modules section.
type Module struct {
	// An MTA internal module name. Names need to be unique within the MTA scope
	Name string
	// a globally unique type ID. Deployment tools will interpret this type ID
	Type string
	// a non-translatable description of this module. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	// A file path which identifies the location of module artifacts.
	Path string `yaml:"path,omitempty"`
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

// Provides List of provided names to which properties (config data) can be attached.
type Provides struct {
	Name       string
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

// Requires list of names either matching a resource name or a name provided by another module within the same MTA.
type Requires struct {
	// an MTA internal name which must match either a provided name, a resource name, or a module name within the same MTA
	Name string `yaml:"name,omitempty"`
	// A group name which shall be use by a deployer to group properties for lookup by a module runtime.
	Group string `yaml:"group,omitempty"`
	// Provided property values can be accessed by "~{<provided-property-name>}". Such expressions can be part of an arbitrary string
	Properties map[string]interface{} `yaml:"properties,omitempty"`
	// Parameters can be used to influence the behavior of tools which interpret this descriptor. Parameters are not made available to requiring modules at runtime
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

// Resource can be anything required to run the application which is not provided by the application itself.
type Resource struct {
	Name string
	// A type of a resource. This type is interpreted by and must be known to the deployer. Resources can be untyped
	Type string
	// A non-translatable description of this resource. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	// Parameters can be used to influence the behavior of tools which interpret this descriptor. Parameters are not made available to requiring modules at runtime
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	// property names and values make up the configuration data which is to be provided to requiring modules at runtime
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

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
