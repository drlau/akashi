package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestResourceValuesGetCombined(t *testing.T) {
	cases := map[string]struct {
		rv       ResourceValues
		expected map[string]interface{}
	}{
		"empty": {
			rv:       ResourceValues{},
			expected: map[string]interface{}{},
		},
		"values only": {
			rv: ResourceValues{
				Values: map[string]interface{}{
					"value1": "hello",
				},
			},
			expected: map[string]interface{}{
				"value1": "hello",
			},
		},
		"computed only": {
			rv: ResourceValues{
				Computed: map[string]interface{}{
					"computed": true,
				},
			},
			expected: map[string]interface{}{
				"computed": true,
			},
		},
		"values and computed": {
			rv: ResourceValues{
				Values: map[string]interface{}{
					"value1": "hello",
				},
				Computed: map[string]interface{}{
					"computed": true,
				},
			},
			expected: map[string]interface{}{
				"value1":   "hello",
				"computed": true,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.rv.GetCombined()
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("(-got, +expected)\n%s", diff)
			}
		})
	}
}
