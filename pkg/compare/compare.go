package compare

import (
	"fmt"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
)

type Resource interface {
	CompareResult(map[string]interface{}) *resource.CompareResult
	Compare(resource.ResourceValues) bool
	Diff(resource.ResourceValues) string
}

func constructNameTypeKey(r plan.ResourcePlan) string {
	return fmt.Sprintf("%s.%s", r.GetType(), r.GetName())
}
