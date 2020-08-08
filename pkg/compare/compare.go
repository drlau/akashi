package compare

import (
	"github.com/drlau/akashi/pkg/plan"
)

type Comparer interface {
	Compare(plan.ResourceChange) bool
	Diff(plan.ResourceChange) string
}
