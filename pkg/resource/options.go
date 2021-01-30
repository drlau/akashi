package resource

import (
	"github.com/drlau/akashi/pkg/ruleset"
)

// TODO: consider moving this to functions
type CompareOptions struct {
	EnforceAll      bool
	IgnoreExtraArgs bool
	IgnoreComputed  bool
	RequireAll      bool
	AutoFail        bool
	IgnoreNoOp      bool
}

func NewCompareOptions(config *ruleset.CompareOptions) *CompareOptions {
	if config == nil {
		return &CompareOptions{}
	}

	return &CompareOptions{
		EnforceAll:      boolFromBoolPointer(config.EnforceAll, false),
		IgnoreExtraArgs: boolFromBoolPointer(config.IgnoreExtraArgs, false),
		IgnoreComputed:  boolFromBoolPointer(config.IgnoreComputed, false),
		RequireAll:      boolFromBoolPointer(config.RequireAll, false),
		AutoFail:        boolFromBoolPointer(config.AutoFail, false),
		IgnoreNoOp:      boolFromBoolPointer(config.IgnoreNoOp, false),
	}
}

func newCompareOptionsWithDefault(config *ruleset.CompareOptions, defaultOpts *CompareOptions) *CompareOptions {
	if defaultOpts == nil {
		return NewCompareOptions(config)
	}

	return &CompareOptions{
		EnforceAll:      boolFromBoolPointer(config.EnforceAll, defaultOpts.EnforceAll),
		IgnoreExtraArgs: boolFromBoolPointer(config.IgnoreExtraArgs, defaultOpts.IgnoreExtraArgs),
		IgnoreComputed:  boolFromBoolPointer(config.IgnoreComputed, defaultOpts.IgnoreComputed),
		RequireAll:      boolFromBoolPointer(config.RequireAll, defaultOpts.RequireAll),
		AutoFail:        boolFromBoolPointer(config.AutoFail, defaultOpts.AutoFail),
		IgnoreNoOp:      boolFromBoolPointer(config.IgnoreNoOp, defaultOpts.IgnoreNoOp),
	}
}

func boolFromBoolPointer(b *bool, failover bool) bool {
	if b != nil {
		return *b
	}
	return failover
}
