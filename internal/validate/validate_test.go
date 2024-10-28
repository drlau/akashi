package validate

import (
	"testing"

	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/google/go-cmp/cmp"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		rs       ruleset.Ruleset
		expected *ValidateResult
	}{
		"valid ruleset without required names": {
			rs: ruleset.Ruleset{
				CreatedResources: &ruleset.CreateDeleteResourceChanges{
					Strict:      false,
					RequireName: false,
					Resources: []ruleset.CreateDeleteResourceChange{
						{
							ResourceIdentifier: ruleset.ResourceIdentifier{
								Type: "google_project_service",
								Name: "api",
							},
						},
						{
							ResourceIdentifier: ruleset.ResourceIdentifier{
								Type: "google_service_account",
							},
						},
					},
				},
			},
			expected: &ValidateResult{},
		},
		"valid ruleset with required names": {
			rs: ruleset.Ruleset{
				CreatedResources: &ruleset.CreateDeleteResourceChanges{
					Strict:      false,
					RequireName: true,
					Resources: []ruleset.CreateDeleteResourceChange{
						{
							ResourceIdentifier: ruleset.ResourceIdentifier{
								Type: "google_project_service",
								Name: "api",
							},
						},
						{
							ResourceIdentifier: ruleset.ResourceIdentifier{
								Type: "google_service_account",
								Name: "my_service_account",
							},
						},
					},
				},
			},
			expected: &ValidateResult{},
		},
		"invalid ruleset": {
			rs: ruleset.Ruleset{
				CreatedResources: &ruleset.CreateDeleteResourceChanges{
					Strict:      false,
					RequireName: true,
					Resources: []ruleset.CreateDeleteResourceChange{
						{
							ResourceIdentifier: ruleset.ResourceIdentifier{
								Type: "google_project_service",
							},
						},
					},
				},
			},
			expected: &ValidateResult{
				InvalidCreatedResources: []*ruleset.ResourceIdentifier{
					{Type: "google_project_service"},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Validate(test.rs)
			if diff := cmp.Diff(res, test.expected); diff != "" {
				t.Errorf("Validate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := map[string]struct {
		res      ValidateResult
		expected bool
	}{
		"no invalid resources": {
			res:      ValidateResult{},
			expected: true,
		},
		"invalid created resources": {
			res: ValidateResult{
				InvalidCreatedResources: []*ruleset.ResourceIdentifier{
					{Type: "fake_created_resource"},
				},
			},
			expected: false,
		},
		"invalid destroyed resources": {
			res: ValidateResult{
				InvalidDestroyedResources: []*ruleset.ResourceIdentifier{
					{Type: "fake_destroy_resource"},
				},
			},
			expected: false,
		},
		"invalid updated resources": {
			res: ValidateResult{
				InvalidUpdatedResources: []*ruleset.ResourceIdentifier{
					{Type: "fake_update_resource"},
				},
			},
			expected: false,
		},
		"multiple invalid resources changes": {
			res: ValidateResult{
				InvalidCreatedResources: []*ruleset.ResourceIdentifier{
					{Type: "fake_created_resource"},
				},
				InvalidUpdatedResources: []*ruleset.ResourceIdentifier{
					{Type: "fake_update_resource"},
				},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.res.IsValid() != test.expected {
				t.Errorf("Expected %t, got %t: %v", test.expected, test.res.IsValid(), test.res)
			}
		})
	}
}
