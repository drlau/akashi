package compare

import (
	"strings"
	"testing"

	comparefakes "github.com/drlau/akashi/pkg/compare/fakes"
	planfakes "github.com/drlau/akashi/pkg/compare/fakes"
	"github.com/drlau/akashi/pkg/plan"
)

func TestUpdateCompare(t *testing.T) {
	cases := map[string]struct {
		comparer     *UpdateComparer
		resourcePlan plan.ResourcePlan
		expected     bool
	}{
		"matching nametype resource with before only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching nametype resource with after only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						After: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching nametype resource with before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
						After: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching nametype resource with passing before and failing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
						After: &comparefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
		"matching nametype resource with failing before and passing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: false,
						},
						After: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
		"matching nametype resource with failing before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: false,
						},
						After: &comparefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
		"matching name resource": {
			comparer: &UpdateComparer{
				NameResources: map[string]updateResource{
					"name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"matching type resource": {
			comparer: &UpdateComparer{
				TypeResources: map[string]updateResource{
					"type": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"prioritizes matching nametype resource": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: true,
						},
					},
				},
				NameResources: map[string]updateResource{
					"name": {
						Before: &comparefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
				TypeResources: map[string]updateResource{
					"type": {
						Before: &comparefakes.FakeResource{
							CompareReturns: false,
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"no matching resource": {
			comparer: &UpdateComparer{},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"no matching resource with strict enabled": {
			comparer: &UpdateComparer{
				Strict: true,
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.comparer.Compare(tc.resourcePlan); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestUpdateDiff(t *testing.T) {
	cases := map[string]struct {
		comparer       *UpdateComparer
		resourcePlan   plan.ResourcePlan
		expected       bool
		expectedOutput []string
	}{
		"matching nametype resource": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching nametype resource with after only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						After: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching nametype resource with before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
						After: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching nametype resource with passing before and failing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
						After: &comparefakes.FakeResource{
							DiffReturns: "failedAfter",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address", "(after)", "failedAfter"},
		},
		"matching nametype resource with failing before and passing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "failedBefore",
						},
						After: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address", "(before)", "failedBefore"},
		},
		"matching nametype resource with failing before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "failedBefore",
						},
						After: &comparefakes.FakeResource{
							DiffReturns: "failedAfter",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address", "(before)", "failedBefore", "address", "(after)", "failedAfter"},
		},
		"matching name resource": {
			comparer: &UpdateComparer{
				NameResources: map[string]updateResource{
					"name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"matching type resource": {
			comparer: &UpdateComparer{
				TypeResources: map[string]updateResource{
					"type": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"prioritizes matching nametype resource": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "",
						},
					},
				},
				NameResources: map[string]updateResource{
					"name": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "failedName",
						},
					},
				},
				TypeResources: map[string]updateResource{
					"type": {
						Before: &comparefakes.FakeResource{
							DiffReturns: "failedType",
						},
					},
				},
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{""},
		},
		"no matching resource": {
			comparer: &UpdateComparer{},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       true,
			expectedOutput: []string{"!", "address", "(no matching rule)"},
		},
		"no matching resource with strict enabled": {
			comparer: &UpdateComparer{
				Strict: true,
			},
			resourcePlan: &planfakes.FakeResourcePlan{
				AddressReturns: "address",
				NameReturns:    "name",
				TypeReturns:    "type",
			},
			expected:       false,
			expectedOutput: []string{"×", "address", "(no matching rule)"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			output, got := tc.comparer.Diff(tc.resourcePlan)
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
