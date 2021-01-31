package compare

import (
	"fmt"
	"strings"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

type UpdateComparer struct {
	Strict bool

	NameResources     map[string]updateResource
	TypeResources     map[string]updateResource
	NameTypeResources map[string]updateResource
}

type updateResource struct {
	Before Resource
	After  Resource
}

func NewUpdateComparer(ruleset ruleset.UpdateResourceChanges) *UpdateComparer {
	defaultOptions := resource.NewCompareOptions(ruleset.Default)
	nameTypeResources := make(map[string]updateResource)
	typeResources := make(map[string]updateResource)
	nameResources := make(map[string]updateResource)

	// Iterate over all the resources
	for _, r := range ruleset.Resources {
		var ur updateResource
		if r.Before != nil {
			ur.Before = resource.NewResourceFromConfig(r.ResourceIdentifier, *r.Before, &r.CompareOptions, defaultOptions)
		}
		if r.After != nil {
			ur.After = resource.NewResourceFromConfig(r.ResourceIdentifier, *r.After, &r.CompareOptions, defaultOptions)
		}

		if r.Name != "" && r.Type != "" {
			// format name and type key
			nameTypeResources[fmt.Sprintf("%s.%s", r.Type, r.Name)] = ur
		} else if r.Name != "" {
			nameResources[r.Name] = ur
		} else if r.Type != "" {
			typeResources[r.Type] = ur
		}
	}
	return &UpdateComparer{
		Strict:            ruleset.Strict,
		NameResources:     nameResources,
		TypeResources:     typeResources,
		NameTypeResources: nameTypeResources,
	}
}

func (c *UpdateComparer) Compare(r plan.ResourcePlan) bool {
	nameType := constructNameTypeKey(r)
	beforeChanges := resource.ResourceValues{
		Values:        r.GetBefore(),
		ChangedValues: r.GetBeforeChangedOnly(),
	}
	afterChanges := resource.ResourceValues{
		Values:        r.GetAfter(),
		ChangedValues: r.GetAfterChangedOnly(),
		Computed:      r.GetComputed(),
	}

	var (
		beforeOk = true
		afterOk  = true
	)
	if ro, ok := c.NameTypeResources[nameType]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.Compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.Compare(afterChanges)
		}
		return beforeOk && afterOk
	} else if ro, ok := c.NameResources[r.GetName()]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.Compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.Compare(afterChanges)
		}
		return beforeOk && afterOk
	} else if ro, ok := c.TypeResources[r.GetType()]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.Compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.Compare(afterChanges)
		}
		return beforeOk && afterOk
	}

	return !c.Strict
}

func (c *UpdateComparer) Diff(r plan.ResourcePlan) (string, bool) {
	nameType := constructNameTypeKey(r)
	// TODO: handle IgnoreNoOp
	beforeChanges := resource.ResourceValues{
		Values:        r.GetBefore(),
		ChangedValues: r.GetBeforeChangedOnly(),
	}
	afterChanges := resource.ResourceValues{
		Values:        r.GetAfter(),
		ChangedValues: r.GetAfterChangedOnly(),
		Computed:      r.GetComputed(),
	}

	var ur updateResource
	if rs, ok := c.NameTypeResources[nameType]; ok {
		ur = rs
	} else if rs, ok := c.NameResources[r.GetName()]; ok {
		ur = rs
	} else if rs, ok := c.TypeResources[r.GetType()]; ok {
		ur = rs
	} else {
		if c.Strict {
			return fmt.Sprintf("%s %s (no matching rule)", utils.Red("×"), r.GetAddress()), false
		}

		return fmt.Sprintf("%s %s (no matching rule)", utils.Yellow("!"), r.GetAddress()), true
	}

	var (
		result strings.Builder
		equal  = true
	)

	if ur.Before != nil {
		diff := ur.Before.Diff(beforeChanges)
		if diff != "" {
			equal = false
			result.WriteString(fmt.Sprintf("%s %s %s\n%s\n", utils.Red("×"), utils.Red(r.GetAddress()), utils.Red("(before)"), diff))
		}
	}

	if ur.After != nil {
		diff := ur.After.Diff(afterChanges)
		if diff != "" {
			equal = false
			result.WriteString(fmt.Sprintf("%s %s %s\n%s\n", utils.Red("×"), utils.Red(r.GetAddress()), utils.Red("(after)"), diff))
		}
	}

	if equal {
		return fmt.Sprintf("%s %s", utils.Green("✓"), r.GetAddress()), true
	}

	return strings.TrimSuffix(result.String(), "\n"), equal
}
