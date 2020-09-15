package resource

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

var (
	emptyMap = map[interface{}]interface{}{}
	emptyStringMap = map[string]interface{}{}
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

	Enforced map[string]interface{}
	Ignored  map[string]interface{}
}

// TODO: consider moving this to functions
type CompareOptions struct {
	EnforceAll      bool
	IgnoreExtraArgs bool
	IgnoreComputed  bool
	RequireAll      bool
	AutoFail        bool
}

func NewResourceFromConfig(c ruleset.ResourceChange) Resource {
	ignored := make(map[string]interface{})

	for _, i := range c.Ignored {
		ignored[i] = true
	}
	return &resource{
		Name:     c.Name,
		Type:     c.Type,
		Enforced: c.Enforced,
		Ignored:  ignored,
	}
}

func (r *resource) CompareResult(values map[string]interface{}) *CompareResult {
	enforcedArgs := make(map[string]interface{})
	failedArgs := make(map[string]interface{})
	ignored := make(map[string]interface{})
	extraArgs := make(map[string]interface{})

	// Passed in the plan's values
	// Iterate over each key/value
	for k, v := range values {
		// If the key is ignored, record and continue
		if _, ok := r.Ignored[k]; ok {
			ignored[k] = true
			continue
		}
		// If the key is enforced...
		if enforced, ok := r.Enforced[k]; ok {
			// YAML parses "key: {}" as a map[interface{}]interface{} which is different from map[string]interface{}
			if mapEnforced, ok := enforced.(map[interface{}]interface{}); ok {
				enforced = convertMapKeysToString(mapEnforced)
			}

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
		Ignored:         ignored,
		Extra:           extraArgs,
		MissingEnforced: setDifference(setDifference(r.Enforced, enforcedArgs), failedArgs),
		MissingIgnored:  setDifference(setDifference(r.Ignored, ignored), failedArgs),
	}
}

func (r *resource) Compare(rv ResourceValues, opts CompareOptions) bool {
	if opts.AutoFail {
		return false
	}
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
	if opts.AutoFail {
		return utils.Red("AutoFail set to true")
	}
	var buf strings.Builder
	values := rv.Values
	if !opts.IgnoreComputed {
		values = rv.GetCombined()
	}
	cmp := r.CompareResult(values)

	if opts.EnforceAll && len(cmp.MissingEnforced) > 0 {
		buf.WriteString(utils.Red("Missing enforced arguments:\n"))
		for arg, _ := range cmp.MissingEnforced {
			buf.WriteString(utils.Red(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if !opts.IgnoreExtraArgs && len(cmp.Extra) != 0 {
		buf.WriteString(utils.Yellow("Extra arguments:\n"))
		for arg, _ := range cmp.Extra {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if opts.RequireAll && (len(cmp.MissingEnforced)+len(cmp.MissingIgnored)) != 0 {
		buf.WriteString(utils.Yellow("Missing enforced and ignored arguments:\n"))
		for arg, _ := range cmp.MissingEnforced {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
		for arg, _ := range cmp.MissingIgnored {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
	}

	if len(cmp.Failed) > 0 {
		buf.WriteString(utils.Red("Failed arguments:\n"))
		for k, v := range cmp.Failed {
			f := v.(FailedArg)

			buf.WriteString(utils.Red(fmt.Sprintf("  - %v\n", k)))
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

// convertMapKeysToString converts map[interface{}]interface{} to map[string]interface{}
func convertMapKeysToString(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = v
	}

	return result
}