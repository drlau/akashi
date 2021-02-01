package resource

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

var (
	emptyMap       = map[interface{}]interface{}{}
	emptyStringMap = map[string]interface{}{}
)

type resource struct {
	Name string
	Type string
	// TODO support Index
	// Index interface{}

	Enforced map[string]ruleset.EnforceChange
	Ignored  map[string]interface{}

	CompareOptions *CompareOptions
}

func NewResourceFromConfig(resourceIdentifier ruleset.ResourceIdentifier, resourceRules ruleset.ResourceRules, resourceOpts *ruleset.CompareOptions, defaultOpts *CompareOptions) *resource {
	ignored := make(map[string]interface{})

	for _, i := range resourceRules.Ignored {
		ignored[i] = true
	}
	return &resource{
		Name:     resourceIdentifier.Name,
		Type:     resourceIdentifier.Type,
		Enforced: resourceRules.Enforced,
		Ignored:  ignored,

		CompareOptions: newCompareOptionsWithDefault(resourceOpts, defaultOpts),
	}
}

func (r *resource) CompareResult(values map[string]interface{}) *CompareResult {
	result := &CompareResult{
		Enforced: make(map[string]interface{}),
		Failed:   make(map[string]interface{}),
		Ignored:  make(map[string]interface{}),
		Extra:    make(map[string]interface{}),
	}

	result.checkValues(r.Enforced, r.Ignored, values, "")

	result.MissingEnforced = setDifference(enforcedSetDifference(make(map[string]interface{}), "", r.Enforced, result.Enforced), result.Failed)
	result.MissingIgnored = setDifference(setDifference(r.Ignored, result.Ignored), result.Failed)

	return result
}

func (r *resource) Compare(rv ResourceValues) bool {
	if r.CompareOptions.AutoFail {
		return false
	}
	values := rv.Values
	if r.CompareOptions.IgnoreNoOp && rv.ChangedValues != nil {
		values = rv.ChangedValues
	} else if !r.CompareOptions.IgnoreComputed {
		values = rv.GetCombined()
	}
	cmp := r.CompareResult(values)

	if r.CompareOptions.EnforceAll && len(cmp.MissingEnforced) > 0 {
		return false
	}
	if !r.CompareOptions.IgnoreExtraArgs && len(cmp.Extra) != 0 {
		return false
	}
	if r.CompareOptions.RequireAll && (len(cmp.MissingEnforced)+len(cmp.MissingIgnored)) != 0 {
		return false
	}

	return len(cmp.Failed) == 0
}

func (r *resource) Diff(rv ResourceValues) string {
	if r.CompareOptions.AutoFail {
		return utils.Red("AutoFail set to true")
	}
	var buf strings.Builder
	values := rv.Values
	if r.CompareOptions.IgnoreNoOp && rv.ChangedValues != nil {
		values = rv.ChangedValues
	} else if !r.CompareOptions.IgnoreComputed {
		values = rv.GetCombined()
	}
	cmp := r.CompareResult(values)

	if r.CompareOptions.EnforceAll && len(cmp.MissingEnforced) > 0 {
		buf.WriteString(utils.Red("Missing enforced arguments:\n"))
		for arg := range cmp.MissingEnforced {
			buf.WriteString(utils.Red(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if !r.CompareOptions.IgnoreExtraArgs && len(cmp.Extra) != 0 {
		buf.WriteString(utils.Yellow("Extra arguments:\n"))
		for arg := range cmp.Extra {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
	}
	if r.CompareOptions.RequireAll && (len(cmp.MissingEnforced)+len(cmp.MissingIgnored)) != 0 {
		buf.WriteString(utils.Yellow("Missing enforced and ignored arguments:\n"))
		for arg := range cmp.MissingEnforced {
			buf.WriteString(utils.Yellow(fmt.Sprintf("  - %v\n", arg)))
		}
		for arg := range cmp.MissingIgnored {
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

func equal(expected, value interface{}) bool {
	// YAML parses "key: {}" as a map[interface{}]interface{} which is different from map[string]interface{}
	if mapExpected, ok := expected.(map[interface{}]interface{}); ok {
		expected = convertMapKeysToString(mapExpected)
	}

	return reflect.DeepEqual(expected, value)
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

// enforcedSetDifference returns elements in A but not in B
// only checks for key equality - ignores values
func enforcedSetDifference(result map[string]interface{}, keyPrefix string, a map[string]ruleset.EnforceChange, b map[string]interface{}) map[string]interface{} {
	for k, v := range a {
		if keyPrefix != "" {
			k = fmt.Sprintf("%s.%s", keyPrefix, k)
		}
		if v.EnforceChange != nil {
			result = enforcedSetDifference(result, k, v.EnforceChange, b)
		} else if _, ok := b[k]; !ok {
			result[k] = v
		}
	}

	return result
}

// convertMapKeysToString converts map[interface{}]interface{} to map[string]interface{}
func convertMapKeysToString(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		if mapV, ok := v.(map[interface{}]interface{}); ok {
			result[fmt.Sprintf("%v", k)] = convertMapKeysToString(mapV)
		} else {
			result[fmt.Sprintf("%v", k)] = v
		}
	}

	return result
}
