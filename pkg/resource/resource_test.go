package resource

import (
	"strings"
	"testing"

	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/google/go-cmp/cmp"
)

func TestResourceCompareResult(t *testing.T) {
	cases := map[string]struct {
		resource *resource
		expected *CompareResult
		values   map[string]interface{}
	}{
		"enforced value matches": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": ruleset.EnforceChange{
						Value: "value",
					},
				},
				Failed:          map[string]interface{}{},
				Ignored:         map[string]interface{}{},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"enforced value does not match": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": "value2",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{},
				Failed: map[string]interface{}{
					"key": FailedArg{
						Expected: "value",
						Actual:   "value2",
					},
				},
				Ignored:         map[string]interface{}{},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"extra value that is ignored": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"ignored": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key":     "value",
				"ignored": "ignored",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": ruleset.EnforceChange{
						Value: "value",
					},
				},
				Failed: map[string]interface{}{},
				Ignored: map[string]interface{}{
					"ignored": true,
				},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"extra value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key":   "value",
				"extra": "sensitive",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": ruleset.EnforceChange{
						Value: "value",
					},
				},
				Failed:  map[string]interface{}{},
				Ignored: map[string]interface{}{},
				Extra: map[string]interface{}{
					"extra": true,
				},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"missing enforced value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
					"second": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": ruleset.EnforceChange{
						Value: "value",
					},
				},
				Failed:  map[string]interface{}{},
				Ignored: map[string]interface{}{},
				Extra:   map[string]interface{}{},
				MissingEnforced: map[string]interface{}{
					"second": ruleset.EnforceChange{
						Value: "value",
					},
				},
				MissingIgnored: map[string]interface{}{},
			},
		},
		"ignored arg match only": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{},
				Failed:   map[string]interface{}{},
				Ignored: map[string]interface{}{
					"key": true,
				},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"ignored arg with extra value": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key":   "value",
				"extra": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{},
				Failed:   map[string]interface{}{},
				Ignored: map[string]interface{}{
					"key": true,
				},
				Extra: map[string]interface{}{
					"extra": true,
				},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
		"missing ignored value": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key":    true,
					"second": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{},
				Failed:   map[string]interface{}{},
				Ignored: map[string]interface{}{
					"key": true,
				},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored: map[string]interface{}{
					"second": true,
				},
			},
		},
		"nested enforced value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						EnforceChange: map[string]ruleset.EnforceChange{
							"nested-key": {
								Value: true,
							},
						},
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: map[string]interface{}{
				"key": map[string]interface{}{
					"nested-key": true,
				},
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key.nested-key": ruleset.EnforceChange{
						Value: true,
					},
				},
				Failed:          map[string]interface{}{},
				Ignored:         map[string]interface{}{},
				Extra:           map[string]interface{}{},
				MissingEnforced: map[string]interface{}{},
				MissingIgnored:  map[string]interface{}{},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.resource.CompareResult(tc.values)
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("(-got, +expected)\n%s", diff)
			}
		})
	}
}

func TestResourceCompare(t *testing.T) {
	cases := map[string]struct {
		resource *resource
		values   ResourceValues
		expected bool
	}{
		"enforced value matches": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: true,
		},
		"enforced value does not match": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value2",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: false,
		},
		"extra value that is ignored": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"ignored": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: true,
		},
		"extra value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: false,
		},
		"extra value with IgnoreExtraArgs enabled": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					IgnoreExtraArgs: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: true,
		},
		"missing enforced value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
					"second": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: true,
		},
		"missing enforced value with EnforceAll": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
					"second": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					EnforceAll: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: false,
		},
		"ignored arg match only": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: true,
		},
		"ignored arg with extra value": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			expected: false,
		},
		"ignored arg with extra value and ignoreExtraArgs enabled": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{
					IgnoreExtraArgs: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			expected: true,
		},
		"values is missing a key from ignored or enforced": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"enforced": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"key":    true,
					"second": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":      "value",
					"enforced": "value",
				},
			},
			expected: true,
		},
		"values is missing a key from ignored or enforced and requireAll is enabled": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"enforced": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"key":    true,
					"second": true,
				},
				CompareOptions: &CompareOptions{
					RequireAll: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":      "value",
					"enforced": "value",
				},
			},
			expected: false,
		},
		"autofail makes result false even if it passes match": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"enforced": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					AutoFail: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: false,
		},
		"autofail with no enforced or ignored": {
			resource: &resource{
				CompareOptions: &CompareOptions{
					AutoFail: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: false,
		},
		"enforced value in matchAny": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						MatchAny: []interface{}{"value1", "value2"},
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value1",
				},
			},
			expected: true,
		},
		"enforced value not in matchAny": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						MatchAny: []interface{}{"value1", "value2"},
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.resource.Compare(tc.values); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

// TODO: should verify order of expected strings
func TestResourceDiff(t *testing.T) {
	cases := map[string]struct {
		resource *resource
		values   ResourceValues
		expected []string
	}{
		"enforced value matches": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{},
		},
		"enforced value does not match": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value2",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{
				"Failed arguments:",
				"- key",
				"+ Expected: value2",
				"- Actual:   value",
			},
		},
		"extra value that is ignored": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"ignored": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: []string{},
		},
		"extra value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: []string{
				"Extra arguments:",
				"ignored",
			},
		},
		"extra value with IgnoreExtraArgs enabled": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					IgnoreExtraArgs: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":     "value",
					"ignored": "ignored",
				},
			},
			expected: []string{},
		},
		"missing enforced value": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
					"second": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{},
		},
		"missing enforced value with EnforceAll": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
					"second": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					EnforceAll: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{
				"Missing enforced arguments:",
				"second",
			},
		},
		"ignored arg match only": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{},
		},
		"ignored arg with extra value": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			expected: []string{
				"Extra arguments:",
				"second",
			},
		},
		"ignored arg with extra value and ignoreExtraArgs enabled": {
			resource: &resource{
				Ignored: map[string]interface{}{
					"key": true,
				},
				CompareOptions: &CompareOptions{
					IgnoreExtraArgs: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			expected: []string{},
		},
		"values is missing a key from ignored or enforced": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"enforced": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"key":    true,
					"second": true,
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":      "value",
					"enforced": "value",
				},
			},
			expected: []string{},
		},
		"values is missing a key from ignored or enforced and requireAll is enabled": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"enforced": {
						Value: "value",
					},
				},
				Ignored: map[string]interface{}{
					"key":    true,
					"second": true,
				},
				CompareOptions: &CompareOptions{
					RequireAll: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":      "value",
					"enforced": "value",
				},
			},
			expected: []string{
				"Missing enforced and ignored arguments:",
				"second",
			},
		},
		"autofail makes result false even if it passes match": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						Value: "value",
					},
				},
				CompareOptions: &CompareOptions{
					AutoFail: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{"AutoFail set to true"},
		},
		"autofail with no enforced or ignored": {
			resource: &resource{
				CompareOptions: &CompareOptions{
					AutoFail: true,
				},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{"AutoFail set to true"},
		},
		"enforced value in matchAny": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						MatchAny: []interface{}{"value1", "value2"},
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value1",
				},
			},
			expected: []string{},
		},
		"enforced value not in matchAny": {
			resource: &resource{
				Enforced: map[string]ruleset.EnforceChange{
					"key": {
						MatchAny: []interface{}{"value1", "value2"},
					},
				},
				CompareOptions: &CompareOptions{},
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{"one of: [value1 value2]"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.resource.Diff(tc.values)
			if len(tc.expected) == 0 && got != "" {
				t.Errorf("Expected no diff but got: \n%v", got)
			}
			for _, s := range tc.expected {
				if !strings.Contains(got, s) {
					t.Errorf("Result string did not contain %v", s)
				}
			}
		})
	}
}

func TestSetDifference(t *testing.T) {
	cases := map[string]struct {
		a        map[string]interface{}
		b        map[string]interface{}
		expected map[string]interface{}
	}{
		"distinct": {
			a: map[string]interface{}{
				"keyA": "valueA",
			},
			b: map[string]interface{}{
				"keyB": "valueB",
			},
			expected: map[string]interface{}{
				"keyA": "valueA",
			},
		},
		"one common element": {
			a: map[string]interface{}{
				"keyA":      "valueA",
				"sharedKey": "sharedValue",
			},
			b: map[string]interface{}{
				"keyB":      "valueB",
				"sharedKey": "sharedValue",
			},
			expected: map[string]interface{}{
				"keyA": "valueA",
			},
		},
		"multiple common element": {
			a: map[string]interface{}{
				"keyA":       "valueA",
				"sharedKey1": "sharedValue1",
				"sharedKey2": "sharedValue2",
			},
			b: map[string]interface{}{
				"keyB":       "valueB",
				"sharedKey1": "sharedValue1",
				"sharedKey2": "sharedValue2",
			},
			expected: map[string]interface{}{
				"keyA": "valueA",
			},
		},
		"identical sets": {
			a: map[string]interface{}{
				"sharedKey1": "sharedValue1",
			},
			b: map[string]interface{}{
				"sharedKey1": "sharedValue1",
			},
			expected: map[string]interface{}{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := setDifference(tc.a, tc.b)
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("(-got, +expected)\n%s", diff)
			}
		})
	}
}
