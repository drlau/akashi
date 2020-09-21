package compare

import (
	"fmt"

	"github.com/drlau/akashi/pkg/resource"
	"github.com/drlau/akashi/pkg/ruleset"
)

type ResourceChange interface {
	IsCreate() bool
	IsDelete() bool
	GetBefore() map[string]interface{}
	GetAfter() map[string]interface{}
	GetComputed() map[string]interface{}
	GetName() string
	GetType() string
	GetAddress() string
}

type resourceWithOpts struct {
	resource resource.Resource
	opts     resource.CompareOptions
}

func newCreateDeleteResourceWithOpts(resourceConfig ruleset.CreateDeleteResourceChange, defaultOptions resource.CompareOptions) resourceWithOpts {
	return newResourceWithOpts(resourceConfig.ResourceIdentifier, resourceConfig.ResourceRules, resourceConfig.CompareOptions, defaultOptions)
}

func newResourceWithOpts(resourceIdentifier ruleset.ResourceIdentifier, resourceRules ruleset.ResourceRules, resourceOpts ruleset.CompareOptions, defaultOptions resource.CompareOptions) resourceWithOpts {
	return resourceWithOpts{
		resource: resource.NewResourceFromConfig(resourceIdentifier, resourceRules),
		opts: resource.CompareOptions{
			EnforceAll:      boolFromBoolPointer(resourceOpts.EnforceAll, defaultOptions.EnforceAll),
			IgnoreExtraArgs: boolFromBoolPointer(resourceOpts.IgnoreExtraArgs, defaultOptions.IgnoreExtraArgs),
			IgnoreComputed:  boolFromBoolPointer(resourceOpts.IgnoreComputed, defaultOptions.IgnoreComputed),
			RequireAll:      boolFromBoolPointer(resourceOpts.RequireAll, defaultOptions.RequireAll),
			AutoFail:        boolFromBoolPointer(resourceOpts.AutoFail, defaultOptions.AutoFail),
			IgnoreNoOp:      boolFromBoolPointer(resourceOpts.IgnoreNoOp, defaultOptions.IgnoreNoOp),
		},
	}
}

func (r resourceWithOpts) compare(rv resource.ResourceValues) bool {
	return r.resource.Compare(rv, r.opts)
}

func (r resourceWithOpts) diff(rv resource.ResourceValues) string {
	return r.resource.Diff(rv, r.opts)
}

func makeDefaultCompareOptions(config *ruleset.CompareOptions) resource.CompareOptions {
	if config == nil {
		return resource.CompareOptions{}
	}
	return resource.CompareOptions{
		EnforceAll:      boolFromBoolPointer(config.EnforceAll, false),
		IgnoreExtraArgs: boolFromBoolPointer(config.IgnoreExtraArgs, false),
		IgnoreComputed:  boolFromBoolPointer(config.IgnoreComputed, false),
		RequireAll:      boolFromBoolPointer(config.RequireAll, false),
		AutoFail:        boolFromBoolPointer(config.AutoFail, false),
		IgnoreNoOp:      boolFromBoolPointer(config.IgnoreNoOp, false),
	}
}

func constructNameTypeKey(r ResourceChange) string {
	return fmt.Sprintf("%s.%s", r.GetType(), r.GetName())
}

func boolFromBoolPointer(b *bool, failover bool) bool {
	if b != nil {
		return *b
	}
	return failover
}
