package validate

import (
	"github.com/drlau/akashi/pkg/ruleset"
)

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
