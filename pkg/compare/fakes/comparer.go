package fakes

import (
	"io"

	"github.com/drlau/akashi/pkg/plan"
)

type FakeComparer struct {
	CompareReturns bool
	DiffReturns    bool
	DiffOutput     []byte
}

func (r *FakeComparer) Compare(rc plan.ResourceChange) bool {
	return r.CompareReturns
}

func (r *FakeComparer) Diff(out io.Writer, rc plan.ResourceChange) bool {
	out.Write(r.DiffOutput)
	return r.DiffReturns
}
