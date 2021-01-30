package compare

import (
	"fmt"

	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/resource"
)

type ResourceChange interface {
	IsCreate() bool
	IsDelete() bool
	IsNoOp() bool
	IsUpdate() bool
	GetBefore() map[string]interface{}
	GetAfter() map[string]interface{}
	GetBeforeChangedOnly() map[string]interface{}
	GetAfterChangedOnly() map[string]interface{}
	GetComputed() map[string]interface{}
	GetName() string
	GetType() string
	GetAddress() string
}

type Resource interface {
	CompareResult(map[string]interface{}) *resource.CompareResult
	Compare(resource.ResourceValues) bool
	Diff(resource.ResourceValues) string
}

type Comparer interface {
	Compare(plan.ResourceChange) bool
	Diff(plan.ResourceChange) (string, bool)
}

func constructNameTypeKey(r ResourceChange) string {
	return fmt.Sprintf("%s.%s", r.GetType(), r.GetName())
}
