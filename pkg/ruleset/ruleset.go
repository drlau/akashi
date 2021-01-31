package ruleset

type Ruleset struct {
	CreatedResources   *CreateDeleteResourceChanges `yaml:"createdResources,omitempty"`
	DestroyedResources *CreateDeleteResourceChanges `yaml:"destroyedResources,omitempty"`
	UpdatedResources   *UpdateResourceChanges       `yaml:"updatedResources,omitempty"`
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
	CompareOptions     `yaml:",inline"`
	ResourceIdentifier `yaml:",inline"`
	ResourceRules      `yaml:",inline"`
}

type UpdateResourceChanges struct {
	// If strict is enabled, all updated resources must match a rule
	Strict bool `yaml:"strict,omitempty"`

	// Default CompareOptions to use for all resources
	Default *CompareOptions `yaml:"default,omitempty"`

	// Resources is a list of resource changes to validate against
	Resources []UpdateResourceChange `yaml:"resources"`
}

type UpdateResourceChange struct {
	CompareOptions     `yaml:",inline"`
	ResourceIdentifier `yaml:",inline"`

	Before *ResourceRules `yaml:"before,omitempty"`
	After  *ResourceRules `yaml:"after,omitempty"`
}

type CompareOptions struct {
	// If enforceAll is enabled, all Enforced must be present
	EnforceAll *bool `yaml:"enforceAll,omitempty"`

	// If ignoreExtraArgs is enabled, extra args not in Enforced or Ignored are ignored
	IgnoreExtraArgs *bool `yaml:"ignoreExtraArgs,omitempty"`

	// If ignoreComputed is enabled, args that result in a computed value are ignored
	// Has no effect on destroyed values
	IgnoreComputed *bool `yaml:"ignoreComputed,omitempty"`

	// If requireAll is enabled, every key in enforced or Ignored must be present
	RequireAll *bool `yaml:"requireAll,omitempty"`

	// If autoFail is enabled, automatically fails before comparison if a matching resource is found
	AutoFail *bool `yaml:"autoFail,omitempty"`

	// If IgnoreNoOp is enabled, skips attributes that have not changed
	// No effect for created or destroyed resource changes
	IgnoreNoOp *bool `yaml:"ignoreNoOp,omitempty"`
}

type ResourceIdentifier struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	// TODO: index
	// Index interface{} `yaml:"index,omitempty"`
}

type ResourceRules struct {
	Enforced map[string]EnforceChange `yaml:"enforced,omitempty"`
	Ignored  []string                 `yaml:"ignored,omitempty"`
}

type EnforceChange struct {
	Value    interface{}   `yaml:"value,omitempty"`
	MatchAny []interface{} `yaml:"matchAny,omitempty"`
}
