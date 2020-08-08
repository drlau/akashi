package resource

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestResourceCompareResult(t *testing.T) {
	cases := map[string]struct {
		resource Resource
		expected *CompareResult
		values   map[string]interface{}
	}{
		"enforced value matches": {
			resource: &resource{
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": "value",
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
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
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
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
				IgnoredArgs: map[string]interface{}{
					"ignored": true,
				},
			},
			values: map[string]interface{}{
				"key":     "value",
				"ignored": "ignored",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": "value",
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
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
			},
			values: map[string]interface{}{
				"key":   "value",
				"extra": "sensitive",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": "value",
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
				EnforcedValues: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			values: map[string]interface{}{
				"key": "value",
			},
			expected: &CompareResult{
				Enforced: map[string]interface{}{
					"key": "value",
				},
				Failed:  map[string]interface{}{},
				Ignored: map[string]interface{}{},
				Extra:   map[string]interface{}{},
				MissingEnforced: map[string]interface{}{
					"second": "value",
				},
				MissingIgnored: map[string]interface{}{},
			},
		},
		"ignored arg match only": {
			resource: &resource{
				IgnoredArgs: map[string]interface{}{
					"key": true,
				},
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
				IgnoredArgs: map[string]interface{}{
					"key": true,
				},
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
				IgnoredArgs: map[string]interface{}{
					"key":    true,
					"second": true,
				},
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
		resource Resource
		opts     CompareOptions
		values   ResourceValues
		expected bool
	}{
		"enforced value matches": {
			resource: &resource{
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
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
				EnforcedValues: map[string]interface{}{
					"key": "value2",
				},
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
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
				IgnoredArgs: map[string]interface{}{
					"ignored": true,
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
		"extra value": {
			resource: &resource{
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
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
				EnforcedValues: map[string]interface{}{
					"key": "value",
				},
			},
			opts: CompareOptions{
				IgnoreExtraArgs: true,
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
				EnforcedValues: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
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
				EnforcedValues: map[string]interface{}{
					"key":    "value",
					"second": "value",
				},
			},
			opts: CompareOptions{
				EnforceAll: true,
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
				IgnoredArgs: map[string]interface{}{
					"key": true,
				},
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
				IgnoredArgs: map[string]interface{}{
					"key": true,
				},
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
				IgnoredArgs: map[string]interface{}{
					"key": true,
				},
			},
			opts: CompareOptions{
				IgnoreExtraArgs: true,
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
				EnforcedValues: map[string]interface{}{
					"enforced": "value",
				},
				IgnoredArgs: map[string]interface{}{
					"key":    true,
					"second": true,
				},
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
				EnforcedValues: map[string]interface{}{
					"enforced": "value",
				},
				IgnoredArgs: map[string]interface{}{
					"key":    true,
					"second": true,
				},
			},
			opts: CompareOptions{
				RequireAll: true,
			},
			values: ResourceValues{
				Values: map[string]interface{}{
					"key":      "value",
					"enforced": "value",
				},
			},
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.resource.Compare(tc.values, tc.opts); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
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
