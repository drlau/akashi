package validate

import (
	"github.com/drlau/akashi/pkg/ruleset"
)

type ValidateResult struct {
	Valid bool
}

func Validate(ruleset ruleset.Ruleset) *ValidateResult {
	return &ValidateResult{}
}
