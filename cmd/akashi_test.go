package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drlau/akashi/internal/compare"
	comparefakes "github.com/drlau/akashi/internal/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
)

func TestRunCompare(t *testing.T) {
	cases := map[string]struct {
		comparers      compare.ComparerSet
		resourceChange []plan.ResourcePlan
		expected       int
	}{
		"create returns false with create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					CompareReturns: false,
				},
			},
			resourceChange: []plan.ResourcePlan{
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
			resourceChange: []plan.ResourcePlan{
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
			resourceChange: []plan.ResourcePlan{
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
			resourceChange: []plan.ResourcePlan{
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
			resourceChange: []plan.ResourcePlan{
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
			if got := runCompare(tc.resourceChange, tc.comparers); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

// TODO: there must be a better way to do pre and post hooks
// ginkgo comes to mind
func TestRunDiff(t *testing.T) {
	cases := map[string]struct {
		comparers      compare.ComparerSet
		resourceChange []plan.ResourcePlan
		preHook        func()
		expected       int
		expectedOutput []string
	}{
		"create returns false with create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			expected:       0,
			expectedOutput: []string{"comparer fail"},
		},
		"create returns true with create resource": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			expected:       0,
			expectedOutput: []string{"comparer ok"},
		},
		"no matching comparer": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			expected:       0,
			expectedOutput: []string{""},
		},
		"no matching comparer with strict enabled": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			preHook: func() {
				strict = true
			},
			expected:       0,
			expectedOutput: []string{"?", "address (no matching comparer)"},
		},
		"create returns true with multiple resources": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourcePlan{
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
			},
			expected:       0,
			expectedOutput: []string{"comparer ok\ncomparer ok"},
		},
		"fails if there is at least 1 failure": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				DestroyComparer: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			expected:       0,
			expectedOutput: []string{"comparer fail", "comparer ok"},
		},
		"returns 1 if there is at least 1 failure and errorOnFail is set": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				DestroyComparer: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			preHook: func() {
				errorOnFail = true
			},
			expected:       1,
			expectedOutput: []string{"comparer fail", "comparer ok"},
		},
		"only outputs failed with failedOnly": {
			comparers: compare.ComparerSet{
				CreateComparer: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				DestroyComparer: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourcePlan{
				&planfakes.FakeResourcePlan{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourcePlan{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			preHook: func() {
				failedOnly = true
			},
			expected:       0,
			expectedOutput: []string{"comparer fail"},
		},
		// TODO: test case to ensure comparers are called correctly(matching type and number of calls)
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// set default vars
			errorOnFail = false
			strict = false
			failedOnly = false

			if tc.preHook != nil {
				tc.preHook()
			}

			var output bytes.Buffer
			if got := runDiff(&output, tc.resourceChange, tc.comparers); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}

			for _, s := range tc.expectedOutput {
				if !strings.Contains(output.String(), s) {
					t.Errorf("Result string did not contain %v", s)
				}
			}
		})
	}
}
