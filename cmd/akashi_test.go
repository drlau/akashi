package cmd

import (
	"testing"

	"github.com/drlau/akashi/pkg/compare"
	"github.com/drlau/akashi/pkg/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
)

func TestRunCompare(t *testing.T) {
	cases := map[string]struct {
		comparers      map[string]compare.Comparer
		resourceChange []plan.ResourceChange
		expected       int
	}{
		"create returns false with create resource": {
			comparers: map[string]compare.Comparer{
				createKey: &fakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourceChange: []plan.ResourceChange{
				&fakes.FakeResourceChange{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 1,
		},
		"create returns true with create resource": {
			comparers: map[string]compare.Comparer{
				createKey: &fakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourceChange: []plan.ResourceChange{
				&fakes.FakeResourceChange{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 0,
		},
		"create returns false with non-create resource": {
			comparers: map[string]compare.Comparer{
				createKey: &fakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourceChange: []plan.ResourceChange{
				&fakes.FakeResourceChange{
					NameReturns: "name",
					TypeReturns: "type",
				},
			},
			expected: 0,
		},
		"create returns true with multiple resources": {
			comparers: map[string]compare.Comparer{
				createKey: &fakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourceChange: []plan.ResourceChange{
				&fakes.FakeResourceChange{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
				&fakes.FakeResourceChange{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 0,
		},
		"fails if there is at least 1 failure": {
			comparers: map[string]compare.Comparer{
				createKey: &fakes.FakeComparer{
					CompareReturns: false,
				},
				destroyKey: &fakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourceChange: []plan.ResourceChange{
				&fakes.FakeResourceChange{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
				&fakes.FakeResourceChange{
					DeleteReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 1,
		},
		// TODO: test case to ensure comparers are called correctly(matching type and number of calls)
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := runCompare(tc.resourceChange, tc.comparers, false); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
