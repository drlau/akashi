package compare

import (
	"strings"
	"testing"

	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
	resourcefakes "github.com/drlau/akashi/pkg/resource/fakes"
)

func TestCreateCompare(t *testing.T) {
	cases := map[string]struct {
		comparer       *CreateComparer
		resourceChange plan.ResourceChange
		expected       bool
	}{
		"matching nametype resource": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching nametype resource returning false": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
		"matching name resource": {
			comparer: &CreateComparer{
				NameResources: map[string]resourceWithOpts{
					"name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching type resource": {
			comparer: &CreateComparer{
				TypeResources: map[string]resourceWithOpts{
					"type": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"prioritizes matching nametype resource": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
				NameResources: map[string]resourceWithOpts{
					"name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
				TypeResources: map[string]resourceWithOpts{
					"type": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"no matching resource": {
			comparer: &CreateComparer{},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"no matching resource with strict enabled": {
			comparer: &CreateComparer{
				Strict: true,
			},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.comparer.Compare(tc.resourceChange); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestCreateDiff(t *testing.T) {
	cases := map[string]struct {
		comparer       *CreateComparer
		resourceChange plan.ResourceChange
		expected       bool
		expectedOutput []string
	}{
		"matching nametype resource": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching nametype resource returning false": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "failed",
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address"},
		},
		"matching name resource": {
			comparer: &CreateComparer{
				NameResources: map[string]resourceWithOpts{
					"name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching type resource": {
			comparer: &CreateComparer{
				TypeResources: map[string]resourceWithOpts{
					"type": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"prioritizes matching nametype resource": {
			comparer: &CreateComparer{
				NameTypeResources: map[string]resourceWithOpts{
					"type.name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
				NameResources: map[string]resourceWithOpts{
					"name": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "failed",
						},
					},
				},
				TypeResources: map[string]resourceWithOpts{
					"type": resourceWithOpts{
						resource: &resourcefakes.FakeResource{
							DiffReturns: "failed",
						},
					},
				},
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"no matching resource": {
			comparer: &CreateComparer{},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{"!", "address (no matching rule)"},
		},
		"no matching resource with strict enabled": {
			comparer: &CreateComparer{
				Strict: true,
			},
			resourceChange: &planfakes.FakeResourceChange{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address (no matching rule)"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			output, got := tc.comparer.Diff(tc.resourceChange)
			if got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
			for _, o := range tc.expectedOutput {
				if !strings.Contains(output, o) {
					t.Errorf("Output %s did not contain expected string %s", output, o)
				}
			}
		})
	}
}
