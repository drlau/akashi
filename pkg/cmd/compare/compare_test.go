package compare

import (
	"testing"

	"github.com/drlau/akashi/internal/compare"
	comparefakes "github.com/drlau/akashi/internal/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
)

func TestRunCompare(t *testing.T) {
	cases := map[string]struct {
		comparers    compare.ComparerSet
		resourcePlan []plan.ResourcePlan
		expected     int
	}{
		"create returns false with create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 1,
		},
		"create returns true with create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 0,
		},
		"create returns false with non-create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					NameReturns: "name",
					TypeReturns: "type",
				},
			},
			expected: 0,
		},
		"create returns true with multiple resources": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
				&planfakes.FakeResourcePlan{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
			},
			expected: 0,
		},
		"fails if there is at least 1 failure": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
				DestroyComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns: true,
					NameReturns:   "name",
					TypeReturns:   "type",
				},
				&planfakes.FakeResourcePlan{
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
			if got := runCompare(tc.resourcePlan, tc.comparers, false); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
