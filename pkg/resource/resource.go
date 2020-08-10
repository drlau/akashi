package resource

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

type Resource interface {
	CompareResult(map[string]interface{}) *CompareResult
	Compare(ResourceValues, CompareOptions) bool
	Diff(ResourceValues, CompareOptions) string
}

type resource struct {
	Name string
	Type string
	// TODO support Index
	// Index interface{}

	EnforcedValues map[string]interface{}
	IgnoredArgs    map[string]interface{}
}

type CompareOptions struct {
	EnforceAll      bool
	IgnoreExtraArgs bool
	IgnoreComputed  bool
	RequireAll      bool
}

func NewResourceFromConfig(c ruleset.ResourceChange) Resource {
	ignoredArgs := make(map[string]interface{})

	for _, i := range c.IgnoredArgs {
		ignoredArgs[i] = true
	}
	return &resource{
		Name:           c.Name,
		Type:           c.Type,
		EnforcedValues: c.EnforcedValues,
		IgnoredArgs:    ignoredArgs,
	}
}

func (r *resource) CompareResult(values map[string]interface{}) *CompareResult {
	enforcedArgs := make(map[string]interface{})
	failedArgs := make(map[string]interface{})
	ignoredArgs := make(map[string]interface{})
	extraArgs := make(map[string]interface{})

	// Passed in the plan's values
	// Iterate over each key/value
	for k, v := range values {
		// If the key is ignored, record and continue
		if _, ok := r.IgnoredArgs[k]; ok {
			ignoredArgs[k] = true
			continue
		}
		// If the key is enforced...
		if enforced, ok := r.EnforcedValues[k]; ok {
			// Verify the value is what is expected
			if !reflect.DeepEqual(v, enforced) {
				// Not equal - record as failed
				failedArgs[k] = FailedArg{
					Expected: enforced,
					Actual:   v,
				}
			} else {
				// equal
				enforcedArgs[k] = enforced
			}
			// key is not enforced or ignored
		} else {
			extraArgs[k] = true
		}
	}

	return &CompareResult{
		Enforced:        enforcedArgs,
		Failed:          failedArgs,
		Ignored:         ignoredArgs,
		Extra:           extraArgs,
		MissingEnforced: setDifference(setDifference(r.EnforcedValues, enforcedArgs), failedArgs),
		MissingIgnored:  setDifference(setDifference(r.IgnoredArgs, ignoredArgs), failedArgs),
	}
}

func (r *resource) Compare(rv ResourceValues, opts CompareOptions) bool {
	values := rv.Values
	if !opts.IgnoreComputed {
		values = rv.GetCombined()
	}
	cmp := r.CompareResult(values)

	if opts.EnforceAll && len(cmp.MissingEnforced) > 0 {
		return false
	}
	if !opts.IgnoreExtraArgs && len(cmp.Extra) != 0 {
		return false
	}
	if opts.RequireAll && (len(cmp.MissingEnforced)+len(cmp.MissingIgnored)) != 0 {
		return false
	}

	return len(cmp.Failed) == 0
}

func (r *resource) Diff(rv ResourceValues, opts CompareOptions) string {
	var buf strings.Builder
	values := rv.Values
	if !opts.IgnoreComputed {
		values = rv.GetCombined()
	}
	cmp := r.CompareResult(values)

	if opts.EnforceAll && len(cmp.MissingEnforced) > 0 {
		buf.WriteString(utils.RedBold("Missing enforced arguments:\n"))
		for arg, _ := range cmp.MissingEnforced {
			buf.WriteString(utils.Red(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if !opts.IgnoreExtraArgs && len(cmp.Extra) != 0 {
		buf.WriteString(utils.YellowBold("Extra arguments:\n"))
		for arg, _ := range cmp.Extra {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if opts.RequireAll && (len(cmp.MissingEnforced)+len(cmp.MissingIgnored)) != 0 {
		buf.WriteString(utils.YellowBold("Missing enforced and ignored arguments:\n"))
		for arg, _ := range cmp.MissingEnforced {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
		for arg, _ := range cmp.MissingIgnored {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
	}

	if len(cmp.Failed) > 0 {
		buf.WriteString(utils.RedBold("Failed arguments:\n"))
		for k, v := range cmp.Failed {
			f := v.(FailedArg)

			buf.WriteString(utils.RedBold(fmt.Sprintf("  - %v\n", k)))
			buf.WriteString(utils.Green(fmt.Sprintf("    + Expected: %v\n", f.Expected)))
			buf.WriteString(utils.Red(fmt.Sprintf("    - Actual:   %v\n", f.Actual)))
		}
	}

	return buf.String()
}

// setDifference returns elements in A but not in B
// only checks for key equality - ignores values
func setDifference(a, b map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range a {
		if _, ok := b[k]; !ok {
			result[k] = v
		}
	}

	return result
}