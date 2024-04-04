package validate

import (
	"fmt"
	"strings"

	"github.com/drlau/akashi/pkg/ruleset"
)

// TODO: currently the only static validation we do is to check if names are
// present for all resources in the ruleset if RequireName is set. Ideally we
// should support more validations, so this struct will need to evolve to
// contain more information about _why_ the resource is invalid.
type ValidateResult struct {
	InvalidCreatedResources   []*ruleset.ResourceIdentifier
	InvalidDestroyedResources []*ruleset.ResourceIdentifier
	InvalidUpdatedResources   []*ruleset.ResourceIdentifier
}

func (r *ValidateResult) fill_defaults() {
	if r.InvalidCreatedResources == nil {
		r.InvalidCreatedResources = make([]*ruleset.ResourceIdentifier, 0)
	}
	if r.InvalidDestroyedResources == nil {
		r.InvalidDestroyedResources = make([]*ruleset.ResourceIdentifier, 0)
	}
	if r.InvalidUpdatedResources == nil {
		r.InvalidUpdatedResources = make([]*ruleset.ResourceIdentifier, 0)
	}
}

func formatResourceIDs(ids []*ruleset.ResourceIdentifier) []string {
	var lines []string
	for _, id := range ids {
		lines = append(lines, fmt.Sprintf("\t- %s", id.String()))
	}
	return lines
}

func (r *ValidateResult) String() string {
	if r.IsValid() {
		return "All resources valid!"
	}

	lines := []string{
		"Found invalid resources in the ruleset:",
		"---------------------------------------",
	}
	if len(r.InvalidCreatedResources) != 0 {
		lines = append(lines, "Invalid Created Resources:")
		lines = append(lines, formatResourceIDs(r.InvalidCreatedResources)...)
	}
	if len(r.InvalidDestroyedResources) != 0 {
		lines = append(lines, "Invalid Destroyed Resources:")
		lines = append(lines, formatResourceIDs(r.InvalidDestroyedResources)...)
	}
	if len(r.InvalidUpdatedResources) != 0 {
		lines = append(lines, "Invalid Updated Resources:")
		lines = append(lines, formatResourceIDs(r.InvalidUpdatedResources)...)
	}
	return strings.Join(lines, "\n")
}

func (r *ValidateResult) IsValid() bool {
	createdValid := len(r.InvalidCreatedResources) == 0
	destroyedValid := len(r.InvalidDestroyedResources) == 0
	updatedValid := len(r.InvalidUpdatedResources) == 0
	return createdValid && destroyedValid && updatedValid
}

func getUnnamedResources[T ruleset.Resource](rs []T) []*ruleset.ResourceIdentifier {
	var res []*ruleset.ResourceIdentifier
	for _, r := range rs {
		id := r.ID()
		if id.Name == "" {
			res = append(res, id)
		}
	}
	return res
}

func Validate(rs ruleset.Ruleset) *ValidateResult {
	res := &ValidateResult{}
	if rs.CreatedResources != nil && rs.CreatedResources.RequireName {
		ids := getUnnamedResources(rs.CreatedResources.Resources)
		res.InvalidCreatedResources = ids
	}
	if rs.DestroyedResources != nil && rs.DestroyedResources.RequireName {
		ids := getUnnamedResources(rs.DestroyedResources.Resources)
		res.InvalidDestroyedResources = ids
	}
	if rs.UpdatedResources != nil && rs.UpdatedResources.RequireName {
		ids := getUnnamedResources(rs.UpdatedResources.Resources)
		res.InvalidUpdatedResources = ids
	}
	return res
}
