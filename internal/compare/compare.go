package compare

import (
	"github.com/drlau/akashi/pkg/compare"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/ruleset"
)

type Comparer interface {
	Compare(plan.ResourcePlan) bool
	Diff(plan.ResourcePlan) (string, bool)
}

type ComparerSet struct {
	CreateComparer  Comparer
	DestroyComparer Comparer
	UpdateComparer  Comparer
}

func NewComparerSet(path string) (ComparerSet, error) {
	result := ComparerSet{}

	rs, err := ruleset.ParseRuleset(path)
	if err != nil {
		return result, err
	}

	if rs.CreatedResources != nil {
		result.CreateComparer = compare.NewCreateComparer(*rs.CreatedResources)
	}
	if rs.DestroyedResources != nil {
		result.DestroyComparer = compare.NewDestroyComparer(*rs.DestroyedResources)
	}
	if rs.UpdatedResources != nil {
		result.UpdateComparer = compare.NewUpdateComparer(*rs.UpdatedResources)
	}

	return result, nil
}
