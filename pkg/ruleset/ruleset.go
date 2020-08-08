package ruleset

type Ruleset struct {
	CreatedResources   *CreateDeleteResourceChanges `yaml:"createdResources,omitempty"`
	DestroyedResources *CreateDeleteResourceChanges `yaml:"destroyedResources,omitempty"`
	// TODO: updated resources
}

type CreateDeleteResourceChanges struct {
	// If strict is enabled, all created or deleted resources must match a rule
	Strict bool `yaml:"strict,omitempty"`

	// Default CompareOptions to use for all resources
	Default *CompareOptions `yaml:"default,omitempty"`

	// Resources is a list of resource changes to validate against
	Resources []CreateDeleteResourceChange `yaml:"resources"`
}

type CreateDeleteResourceChange struct {
	CompareOptions `yaml:",inline"`
	ResourceChange `yaml:",inline"`
}

type CompareOptions struct {
	// If enforceAll is enabled, all EnforcedValues must be present
	EnforceAll *bool `yaml:"enforceAll,omitempty"`

	// If ignoreExtraArgs is enabled, extra args not in Enforced or IgnoredArgs are ignored
	IgnoreExtraArgs *bool `yaml:"ignoreExtraArgs,omitempty"`

	// If ignoreComputed is enabled, args that result in a computed value are ignored
	// Has no effect on destroyed values
	IgnoreComputed *bool `yaml:"ignoreComputed,omitempty"`

	// If requireAll is enabled, every key in enforcedValues or IgnoredArgs must be present
	RequireAll *bool `yaml:"requireAll,omitempty"`
}

type ResourceChange struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	// TODO: index
	// Index interface{} `yaml:"index,omitempty"`

	EnforcedValues map[string]interface{} `yaml:"enforcedValues,omitempty"`
	IgnoredArgs    []string               `yaml:"ignoredArgs,omitempty"`
}
