package compare

import (
	"fmt"
	"io"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

// TODO: Create and destroy are nearly identical
// depending on how updated resources comparison is implemented, move common logic to internal struct
type CreateComparer struct {
	Strict bool

	NameResources     map[string]resourceWithOpts
	TypeResources     map[string]resourceWithOpts
	NameTypeResources map[string]resourceWithOpts
}

func NewCreateComparer(ruleset ruleset.CreateDeleteResourceChanges) *CreateComparer {
	defaultOptions := makeDefaultCompareOptions(ruleset.Default)
	nameTypeResources := make(map[string]resourceWithOpts)
	typeResources := make(map[string]resourceWithOpts)
	nameResources := make(map[string]resourceWithOpts)

	// Iterate over all the resources
	for _, r := range ruleset.Resources {
		if r.Name != "" && r.Type != "" {
			// format name and type key
			// construct Resource and add to map
			nameTypeResources[fmt.Sprintf("%s.%s", r.Type, r.Name)] = newResourceWithOpts(r, defaultOptions)
		} else if r.Name != "" {
			// construct resource and add to name map
			nameResources[r.Name] = newResourceWithOpts(r, defaultOptions)
		} else if r.Type != "" {
			// construct type and add to type map
			typeResources[r.Type] = newResourceWithOpts(r, defaultOptions)
		}
	}
	return &CreateComparer{
		Strict:            ruleset.Strict,
		NameResources:     nameResources,
		TypeResources:     typeResources,
		NameTypeResources: nameTypeResources,
	}
}

func (c *CreateComparer) Compare(r plan.ResourceChange) bool {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
	}

	if ro, ok := c.NameTypeResources[nameType]; ok {
		return ro.compare(changes)
	} else if ro, ok := c.NameResources[r.GetName()]; ok {
		return ro.compare(changes)
	} else if ro, ok := c.TypeResources[r.GetType()]; ok {
		return ro.compare(changes)
	}

	return !c.Strict
}

func (c *CreateComparer) Diff(out io.Writer, r plan.ResourceChange) bool {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values:   r.GetAfter(),
		Computed: r.GetComputed(),
	}

	var ro resourceWithOpts
	if rs, ok := c.NameTypeResources[nameType]; ok {
		ro = rs
	} else if rs, ok := c.NameResources[r.GetName()]; ok {
		ro = rs
	} else if rs, ok := c.TypeResources[r.GetType()]; ok {
		ro = rs
	} else {
		if c.Strict {
			fmt.Fprintln(out, utils.Red(fmt.Sprintf("× %s (no matching rule)", r.GetAddress())))
			return false
		}

		fmt.Fprintln(out, utils.Yellow(fmt.Sprintf("? %s (no matching rule)", r.GetAddress())))
		return true
	}

	diff := ro.diff(changes)
	if diff != "" {
		fmt.Fprintln(out, utils.Red(fmt.Sprintf("× %s", r.GetAddress())))
		fmt.Fprintln(out, diff)

		return false
	}

	return true
}
