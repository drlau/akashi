package fakes

import (
	"github.com/drlau/akashi/pkg/resource"
)

type FakeResource struct {
	CompareResultReturns *resource.CompareResult
	CompareReturns       bool
	DiffReturns          string
}

func (r *FakeResource) CompareResult(values map[string]interface{}) *resource.CompareResult {
	return r.CompareResultReturns
}

func (r *FakeResource) Compare(rv resource.ResourceValues) bool {
	return r.CompareReturns
}

func (r *FakeResource) Diff(rv resource.ResourceValues) string {
	return r.DiffReturns
}
