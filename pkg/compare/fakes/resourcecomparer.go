package fakes

import (
	"github.com/drlau/akashi/pkg/resource"
)

type FakeResourceComparer struct {
	CompareReturns bool
	DiffReturns    string
}

func (r *FakeResourceComparer) Compare(rv resource.ResourceValues) bool {
	return r.CompareReturns
}

func (r *FakeResourceComparer) Diff(rv resource.ResourceValues) string {
	return r.DiffReturns
}
