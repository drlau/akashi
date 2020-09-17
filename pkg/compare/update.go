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
	Before *resourceWithOpts
	After  *resourceWithOpts
}

func NewUpdateComparer(ruleset ruleset.UpdateResourceChanges) *UpdateComparer {
	defaultOptions := makeDefaultCompareOptions(ruleset.Default)
	nameTypeResources := make(map[string]updateResource)
	typeResources := make(map[string]updateResource)
	nameResources := make(map[string]updateResource)

	// Iterate over all the resources
	for _, r := range ruleset.Resources {
		var ur updateResource
		if r.Before != nil {
			ro := newResourceWithOpts(r.ResourceIdentifier, *r.Before, r.CompareOptions, defaultOptions)
			ur.Before = &ro
		}
		if r.After != nil {
			ro := newResourceWithOpts(r.ResourceIdentifier, *r.After, r.CompareOptions, defaultOptions)
			ur.After = &ro
		}

		if r.Name != "" && r.Type != "" {
			// format name and type key
			// construct Resource and add to map
			nameTypeResources[fmt.Sprintf("%s.%s", r.Type, r.Name)] = ur
		} else if r.Name != "" {
			// construct resource and add to name map
			nameResources[r.Name] = ur
		} else if r.Type != "" {
			// construct type and add to type map
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

func (c *UpdateComparer) Compare(r plan.ResourceChange) bool {
	nameType := constructNameTypeKey(r)
	beforeChanges := resource.ResourceValues{
		Values: r.GetBefore(),
	}
	afterChanges := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
	}

	var (
		beforeOk = true
		afterOk  = true
	)
	if ro, ok := c.NameTypeResources[nameType]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.compare(afterChanges)
		}
		return beforeOk && afterOk
	} else if ro, ok := c.NameResources[r.GetName()]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.compare(afterChanges)
		}
		return beforeOk && afterOk
	} else if ro, ok := c.TypeResources[r.GetType()]; ok {
		if ro.Before != nil {
			beforeOk = ro.Before.compare(beforeChanges)
		}
		if ro.After != nil {
			afterOk = ro.After.compare(afterChanges)
		}
		return beforeOk && afterOk
	}

	return !c.Strict
}

func (c *UpdateComparer) Diff(r plan.ResourceChange) (string, bool) {
	nameType := constructNameTypeKey(r)
	// TODO: handle IgnoreNoOp
	beforeChanges := resource.ResourceValues{
		Values: r.GetBefore(),
	}
	afterChanges := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
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
		diff := ur.Before.diff(beforeChanges)
		if diff != "" {
			equal = false
			result.WriteString(fmt.Sprintf("%s %s %s\n%s\n", utils.Red("×"), utils.Red(r.GetAddress()), utils.Red("(before)"), diff))
		}
	}

	if ur.After != nil {
		diff := ur.After.diff(afterChanges)
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
