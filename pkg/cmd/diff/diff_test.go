package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drlau/akashi/pkg/compare"
	comparefakes "github.com/drlau/akashi/pkg/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
)

func TestRunDiff(t *testing.T) {
	cases := map[string]struct {
		comparers      map[string]compare.Comparer
		resourceChange []plan.ResourceChange
		opts           *DiffOptions
		expected       int
		expectedOutput []string
	}{
		"create returns false with create resource": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &DiffOptions{},
			expected:       0,
			expectedOutput: []string{"comparer fail"},
		},
		"create returns true with create resource": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &DiffOptions{},
			expected:       0,
			expectedOutput: []string{"comparer ok"},
		},
		"no matching comparer": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &DiffOptions{},
			expected:       0,
			expectedOutput: []string{""},
		},
		"no matching comparer with strict enabled": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					AddressReturns: "address",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts: &DiffOptions{
				Strict: true,
			},
			expected:       0,
			expectedOutput: []string{"?", "address (no matching comparer)"},
		},
		"create returns true with multiple resources": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &DiffOptions{},
			expected:       0,
			expectedOutput: []string{"comparer ok\ncomparer ok"},
		},
		"fails if there is at least 1 failure": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				destroyKey: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourceChange{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts:           &DiffOptions{},
			expected:       0,
			expectedOutput: []string{"comparer fail", "comparer ok"},
		},
		"returns 1 if there is at least 1 failure and errorOnFail is set": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				destroyKey: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourceChange{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts: &DiffOptions{
				ErrorOnFail: true,
			},
			expected:       1,
			expectedOutput: []string{"comparer fail", "comparer ok"},
		},
		"only outputs failed with failedOnly": {
			comparers: map[string]compare.Comparer{
				createKey: &comparefakes.FakeComparer{
					DiffReturns: false,
					DiffOutput:  "comparer fail",
				},
				destroyKey: &comparefakes.FakeComparer{
					DiffReturns: true,
					DiffOutput:  "comparer ok",
				},
			},
			resourceChange: []plan.ResourceChange{
				&planfakes.FakeResourceChange{
					CreateReturns:  true,
					AddressReturns: "address1",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
				&planfakes.FakeResourceChange{
					DeleteReturns:  true,
					AddressReturns: "address2",
					NameReturns:    "name",
					TypeReturns:    "type",
				},
			},
			opts: &DiffOptions{
				FailedOnly: true,
			},
			expected:       0,
			expectedOutput: []string{"comparer fail"},
		},
		// TODO: test case to ensure comparers are called correctly(matching type and number of calls)
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			if got := runDiff(&output, tc.resourceChange, tc.comparers, tc.opts); got != tc.expected {
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
