package compare

import (
	"fmt"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

// TODO: Create and destroy are nearly identical
// depending on how updated resources comparison is implemented, move common logic to internal struct
type CreateComparer struct {
	Strict bool

	NameResources     map[string]Resource
	TypeResources     map[string]Resource
	NameTypeResources map[string]Resource
}

func NewCreateComparer(ruleset ruleset.CreateDeleteResourceChanges) *CreateComparer {
	defaultOptions := resource.NewCompareOptions(ruleset.Default)
	nameTypeResources := make(map[string]Resource)
	typeResources := make(map[string]Resource)
	nameResources := make(map[string]Resource)

	// Iterate over all the resources
	for _, r := range ruleset.Resources {
		res := resource.NewResourceFromConfig(r.ResourceIdentifier, r.ResourceRules, &r.CompareOptions, defaultOptions)
		if r.Name != "" && r.Type != "" {
			// format name and type key
			nameTypeResources[fmt.Sprintf("%s.%s", r.Type, r.Name)] = res
		} else if r.Name != "" {
			nameResources[r.Name] = res
		} else if r.Type != "" {
			typeResources[r.Type] = res
		}
	}
	return &CreateComparer{
		Strict:            ruleset.Strict,
		NameResources:     nameResources,
		TypeResources:     typeResources,
		NameTypeResources: nameTypeResources,
	}
}

func (c *CreateComparer) Compare(r plan.ResourcePlan) bool {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
	}

	if ro, ok := c.NameTypeResources[nameType]; ok {
		return ro.Compare(changes)
	} else if ro, ok := c.NameResources[r.GetName()]; ok {
		return ro.Compare(changes)
	} else if ro, ok := c.TypeResources[r.GetType()]; ok {
		return ro.Compare(changes)
	}

	return !c.Strict
}

func (c *CreateComparer) Diff(r plan.ResourcePlan) (string, bool) {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
	}

	var ro Resource
	if rs, ok := c.NameTypeResources[nameType]; ok {
		ro = rs
	} else if rs, ok := c.NameResources[r.GetName()]; ok {
		ro = rs
	} else if rs, ok := c.TypeResources[r.GetType()]; ok {
		ro = rs
	} else {
		if c.Strict {
			return fmt.Sprintf("%s %s (no matching rule)", utils.Red("×"), r.GetAddress()), false
		}

		return fmt.Sprintf("%s %s (no matching rule)", utils.Yellow("!"), r.GetAddress()), true
	}

	diff := ro.Diff(changes)
	if diff != "" {
		return fmt.Sprintf("%s %s\n%s", utils.Red("×"), utils.Red(r.GetAddress()), diff), false
	}

	return fmt.Sprintf("%s %s", utils.Green("✓"), r.GetAddress()), true
}
