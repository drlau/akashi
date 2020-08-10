package fakes

import (
	"github.com/drlau/akashi/pkg/plan"
)

type FakeComparer struct {
	CompareReturns bool
	DiffReturns    bool
	DiffOutput     string
}

func (r *FakeComparer) Compare(rc plan.ResourceChange) bool {
	return r.CompareReturns
}

func (r *FakeComparer) Diff(rc plan.ResourceChange) (string, bool) {
	return r.DiffOutput, r.DiffReturns
}
