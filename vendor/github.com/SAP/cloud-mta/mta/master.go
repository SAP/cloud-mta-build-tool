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
	// Module type declarations
	ModuleTypes []*ModuleTypes `yaml:"module-types,omitempty"`
	// Resource declarations. Resources can be anything required to run the application which is not provided by the application itself
	Resources []*Resource `yaml:"resources,omitempty"`
	// Resource type declarations
	ResourceTypes []*ResourceTypes `yaml:"resource-types,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	// Experimental - use for pre/post hook
	BuildParams *ProjectBuild `yaml:"build-parameters,omitempty"`
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
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
	// THE 'includes' ELEMENT IS ONLY RELEVANT FOR DEVELOPMENT DESCRIPTORS (PRIO TO BUILD), NOT FOR DEPLOYMENT DESCRIPTORS!
	Includes []Includes `yaml:"includes,omitempty"`
	// list of names either matching a resource name or a name provided by another module within the same MTA
	Requires []Requires `yaml:"requires,omitempty"`
	// List of provided names (MTA internal)to which properties (= configuration data) can be attached
	Provides []Provides `yaml:"provides,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor. Parameters are not made available to the module at runtime
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	ParametersMetaData map[string]interface{} `yaml:"parameters-metadata,omitempty"`
	// Build-parameters are specifically steering the behavior of build tools.
	BuildParams map[string]interface{} `yaml:"build-parameters,omitempty"`
	// A list containing the names of the modules that must be deployed prior to this one.
	DeployedAfter interface{} `yaml:"deployed-after,omitempty"`
}

// ModuleTypes module types declarations
type ModuleTypes struct {
	// An MTA internal name of the module type. Can be specified in the 'type' element of modules
	Name string `yaml:"name,omitempty"`
	// The name of the extended type. Can be another resource type defined in this descriptor or one of the default types supported by the deployer
	Extends string `yaml:"extends,omitempty"`
	// Properties inherited by all resources of this type
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
	// Parameters inherited by all resources of this type
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	ParametersMetaData map[string]interface{} `yaml:"parameters-metadata,omitempty"`
}

// Provides List of provided names to which properties (config data) can be attached.
type Provides struct {
	Name string
	// Indicates, that the provided properties shall be made publicly available by the deployer
	Public bool `yaml:"public,omitempty"`
	// property names and values make up the configuration data which is to be provided to requiring modules at runtime
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
}

// Requires list of names either matching a resource name or a name provided by another module within the same MTA.
type Requires struct {
	// an MTA internal name which must match either a provided name, a resource name, or a module name within the same MTA
	Name string `yaml:"name,omitempty"`
	// A group name which shall be use by a deployer to group properties for lookup by a module runtime.
	Group string `yaml:"group,omitempty"`
	// All required and found configuration data sets will be assembled into a JSON array and provided to the module by the lookup name as specified by the value of 'list'
	List string `yaml:"list,omitempty"`
	// Provided property values can be accessed by "~{<provided-property-name>}". Such expressions can be part of an arbitrary string
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
	// Parameters can be used to influence the behavior of tools which interpret this descriptor. Parameters are not made available to requiring modules at runtime
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	ParametersMetaData map[string]interface{} `yaml:"parameters-metadata,omitempty"`
	// THE 'includes' ELEMENT IS ONLY RELEVANT FOR DEVELOPMENT DESCRIPTORS (PRIO TO BUILD), NOT FOR DEPLOYMENT DESCRIPTORS!
	Includes []Includes `yaml:"includes,omitempty"`
}

// Resource can be anything required to run the application which is not provided by the application itself.
type Resource struct {
	Name string
	// A type of a resource. This type is interpreted by and must be known to the deployer. Resources can be untyped
	Type string
	// A non-translatable description of this resource. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	// Parameters can be used to influence the behavior of tools which interpret this descriptor. Parameters are not made available to requiring modules at runtime
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	ParametersMetaData map[string]interface{} `yaml:"parameters-metadata,omitempty"`
	// property names and values make up the configuration data which is to be provided to requiring modules at runtime
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
	// THE 'includes' ELEMENT IS ONLY RELEVANT FOR DEVELOPMENT DESCRIPTORS (PRIO TO BUILD), NOT FOR DEPLOYMENT DESCRIPTORS!
	Includes []Includes `yaml:"includes,omitempty"`
	// A resource can be declared to be optional, if the MTA can compensate for its non-existence
	Optional bool `yaml:"optional,omitempty"`
	// If a resource is declared to be active, it is allocated and bound according to declared requirements
	Active bool `yaml:"active,omitempty"`
	// list of names either matching a resource name or a name provided by another module within the same MTA
	Requires []Requires `yaml:"requires,omitempty"`
}

// ResourceTypes resources type declarations
type ResourceTypes struct {
	// An MTA internal name of the module type. Can be specified in the 'type' element of modules
	Name string `yaml:"name,omitempty"`
	// The name of the extended type. Can be another resource type defined in this descriptor or one of the default types supported by the deployer
	Extends string `yaml:"extends,omitempty"`
	// Properties inherited by all resources of this type
	Properties         map[string]interface{} `yaml:"properties,omitempty"`
	PropertiesMetaData map[string]interface{} `yaml:"properties-metadata,omitempty"`
	// Parameters inherited by all resources of this type
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	ParametersMetaData map[string]interface{} `yaml:"parameters-metadata,omitempty"`
}

// Includes The 'includes' element only relevant for development descriptor, not for deployment descriptor
type Includes struct {
	// A name of an include section. This name will be used by a builder to generate a parameter section in the deployment descriptor
	Name string `yaml:"name,omitempty"`
	// A path pointing to a file which contains a map of parameters, either in JSON or in YAML format.
	Path string `yaml:"path,omitempty"`
}

// ProjectBuild - experimental use for pre/post build hook
type ProjectBuild struct {
	BeforeAll struct {
		Builders Builders `yaml:"builders,omitempty"`
	} `yaml:"before-all,omitempty"`
	AfterAll struct {
		Builders Builders `yaml:"builders,omitempty"`
	} `yaml:"after-all,omitempty"`
}

// Builders - generic builder
type Builders []struct {
	Builder           string `yaml:"builder,omitempty"`
	Timeout           string `yaml:"timeout,omitempty"`
	BuildArtifactName string `yaml:"build-artifact-name,omitempty"`
	Options           struct {
		Execute []string `yaml:"execute,omitempty"`
	} `yaml:"options,omitempty"`
}
