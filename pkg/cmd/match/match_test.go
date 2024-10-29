package match

import (
	"bytes"
	"testing"

	"github.com/drlau/akashi/internal/compare"
	comparefakes "github.com/drlau/akashi/internal/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
)

func TestRunMatch(t *testing.T) {
	cases := map[string]struct {
		comparers      compare.ComparerSet
		resourcePlan   []plan.ResourcePlan
		opts           *MatchOptions
		expectedOutput string
	}{
		"outputs all matching resources from the ruleset": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
				DestroyComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address3",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &MatchOptions{Separator: "\n"},
			expectedOutput: "address1\naddress2\n",
		},
		"outputs all non-matching resources from ruleset when inverted": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
				DestroyComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address3",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &MatchOptions{Separator: "\n", Invert: true},
			expectedOutput: "address3\n",
		},
		"outputs all resources using custom separator": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: true,
				},
				DestroyComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourcePlan: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address3",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &MatchOptions{Separator: ","},
			expectedOutput: "address1,address2\n",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			runMatch(&output, tc.resourcePlan, tc.comparers, tc.opts)

			if output.String() != tc.expectedOutput {
				t.Errorf("Expected: %q\nGot: %q\n", tc.expectedOutput, output.String())
			}
		})
	}
}
