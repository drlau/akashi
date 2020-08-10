package compare

import (
	"fmt"
	"io"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

type DestroyComparer struct {
	Strict bool

	NameResources     map[string]resourceWithOpts
	TypeResources     map[string]resourceWithOpts
	NameTypeResources map[string]resourceWithOpts
}

func NewDestroyComparer(ruleset ruleset.CreateDeleteResourceChanges) *DestroyComparer {
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
	return &DestroyComparer{
		Strict:            ruleset.Strict,
		NameResources:     nameResources,
		TypeResources:     typeResources,
		NameTypeResources: nameTypeResources,
	}
}

func (c *DestroyComparer) Compare(r plan.ResourceChange) bool {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values: r.GetBefore(),
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

func (c *DestroyComparer) Diff(out io.Writer, r plan.ResourceChange) bool {
	nameType := constructNameTypeKey(r)
	changes := resource.ResourceValues{
		Values: r.GetBefore(),
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
