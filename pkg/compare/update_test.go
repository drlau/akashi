package compare

import (
	"strings"
	"testing"

	"github.com/drlau/akashi/pkg/plan"
	planfakes "github.com/drlau/akashi/pkg/plan/fakes"
	resourcefakes "github.com/drlau/akashi/pkg/resource/fakes"
)

func TestUpdateCompare(t *testing.T) {
	cases := map[string]struct {
		comparer       *UpdateComparer
		resourceChange plan.ResourceChange
		expected       bool
	}{
		"matching nametype resource with before only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
		"matching nametype resource with after only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
		"matching nametype resource with before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
		"matching nametype resource with passing before and failing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
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
		"matching nametype resource with failing before and passing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
		"matching nametype resource with failing before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
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
			comparer: &UpdateComparer{
				NameResources: map[string]updateResource{
					"name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
			comparer: &UpdateComparer{
				TypeResources: map[string]updateResource{
					"type": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
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
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: true,
							},
						},
					},
				},
				NameResources: map[string]updateResource{
					"name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
						},
					},
				},
				TypeResources: map[string]updateResource{
					"type": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								CompareReturns: false,
							},
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
			comparer: &UpdateComparer{},
			resourceChange: &planfakes.FakeResourceChange{
				NameReturns: "name",
				TypeReturns: "type",
			},
			expected: true,
		},
		"no matching resource with strict enabled": {
			comparer: &UpdateComparer{
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

func TestUpdateDiff(t *testing.T) {
	cases := map[string]struct {
		comparer       *UpdateComparer
		resourceChange plan.ResourceChange
		expected       bool
		expectedOutput []string
	}{
		"matching nametype resource": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
		"matching nametype resource with after only": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
		"matching nametype resource with before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
		"matching nametype resource with passing before and failing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedAfter",
							},
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
			expectedOutput: []string{"×", "address", "(after)", "failedAfter"},
		},
		"matching nametype resource with failing before and passing after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedBefore",
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
			expectedOutput: []string{"×", "address", "(before)", "failedBefore"},
		},
		"matching nametype resource with failing before and after": {
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedBefore",
							},
						},
						After: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedAfter",
							},
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
			expectedOutput: []string{"×", "address", "(before)", "failedBefore", "address", "(after)", "failedAfter"},
		},
		"matching name resource": {
			comparer: &UpdateComparer{
				NameResources: map[string]updateResource{
					"name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
			comparer: &UpdateComparer{
				TypeResources: map[string]updateResource{
					"type": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
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
			comparer: &UpdateComparer{
				NameTypeResources: map[string]updateResource{
					"type.name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "",
							},
						},
					},
				},
				NameResources: map[string]updateResource{
					"name": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedName",
							},
						},
					},
				},
				TypeResources: map[string]updateResource{
					"type": updateResource{
						Before: &resourceWithOpts{
							resource: &resourcefakes.FakeResource{
								DiffReturns: "failedType",
							},
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
			comparer: &UpdateComparer{},
			resourceChange: &planfakes.FakeResourceChange{
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
			resourceChange: &planfakes.FakeResourceChange{
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
