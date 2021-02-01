package resource

import (
	"fmt"

	"github.com/drlau/akashi/pkg/ruleset"
)

type CompareResult struct {
	// args that had a matching EnforcedValue and are equal
	Enforced map[string]interface{}

	// args that had a matching EnforcedValue but were not equal
	Failed map[string]interface{}

	// args that had a matching Ignored entry
	Ignored map[string]interface{}

	// args that did not have an enforcedValue and were not ignored
	Extra map[string]interface{}

	// args defined in EnforcedValue but were not present
	MissingEnforced map[string]interface{}

	// args defined in Ignored but were not present
	MissingIgnored map[string]interface{}
}

func (cr *CompareResult) checkValues(enforced map[string]ruleset.EnforceChange, ignored map[string]interface{}, values map[string]interface{}, keyPrefix string) {
	// Passed in the plan's values
	// Iterate over each key/value
	for k, v := range values {
		key := k
		if keyPrefix != "" {
			key = fmt.Sprintf("%s.%s", keyPrefix, k)
		}
		// If the key is ignored, record and continue
		if _, ok := ignored[k]; ok {
			// TODO: handle nested ignore
			cr.Ignored[k] = true
			continue
		}
		// If the key is enforced...
		if enforced, ok := enforced[k]; ok {
			switch {
			case enforced.Value != nil:
				// Verify the value is what is expected
				if !equal(enforced.Value, v) {
					// Not equal - record as failed
					cr.Failed[key] = FailedArg{
						Expected: enforced.Value,
						Actual:   v,
					}
				} else {
					// equal
					cr.Enforced[key] = enforced
				}
			case enforced.MatchAny != nil:
				found := false
				for _, val := range enforced.MatchAny {
					// Verify the value is what is expected
					if equal(val, v) {
						// equal
						cr.Enforced[key] = enforced
						found = true
						break
					}
				}
				if !found {
					cr.Failed[key] = FailedArg{
						Expected: fmt.Sprintf("one of: %v", enforced.MatchAny),
						Actual:   v,
						MatchAny: true,
					}
				}
			case enforced.EnforceChange != nil:
				casted, ok := values[k].(map[string]interface{})
				if ok {
					cr.checkValues(enforced.EnforceChange, ignored, casted, k)
				} else {
					// failed
					fmt.Println("failed to cast - failed enforced")
				}
			default:
				// TODO: Tests that key exists and that's it - intended?
			}
		} else {
			cr.Extra[key] = true
		}
	}
}

func (cr *CompareResult) GetEnforced() map[string]interface{} {
	return cr.Enforced
}

func (cr *CompareResult) GetFailed() map[string]interface{} {
	return cr.Failed
}

func (cr *CompareResult) GetIgnored() map[string]interface{} {
	return cr.Ignored
}

func (cr *CompareResult) GetExtra() map[string]interface{} {
	return cr.Extra
}

func (cr *CompareResult) GetMissingEnforced() map[string]interface{} {
	return cr.MissingEnforced
}

func (cr *CompareResult) GetMissingIgnored() map[string]interface{} {
	return cr.MissingIgnored
}

type FailedArg struct {
	Expected interface{}
	Actual   interface{}
	MatchAny bool
}
