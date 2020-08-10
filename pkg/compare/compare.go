package compare

import (
	"io"

	"github.com/drlau/akashi/pkg/plan"
)

type Comparer interface {
	Compare(plan.ResourceChange) bool
	Diff(io.Writer, plan.ResourceChange) bool
}
