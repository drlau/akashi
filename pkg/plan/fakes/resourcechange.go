package fakes

type FakeResourceChange struct {
	AddressReturns  string
	NameReturns     string
	TypeReturns     string
	CreateReturns   bool
	DeleteReturns   bool
	UpdateReturns   bool
	BeforeReturns   map[string]interface{}
	AfterReturns    map[string]interface{}
	ComputedReturns map[string]interface{}
}

func (r *FakeResourceChange) GetAddress() string {
	return r.AddressReturns
}

func (r *FakeResourceChange) GetName() string {
	return r.NameReturns
}

func (r *FakeResourceChange) GetType() string {
	return r.TypeReturns
}

func (r *FakeResourceChange) IsCreate() bool {
	return r.CreateReturns
}

func (r *FakeResourceChange) IsDelete() bool {
	return r.DeleteReturns
}

func (r *FakeResourceChange) IsUpdate() bool {
	return r.UpdateReturns
}

func (r *FakeResourceChange) GetBefore() map[string]interface{} {
	return r.BeforeReturns
}

func (r *FakeResourceChange) GetAfter() map[string]interface{} {
	return r.AfterReturns
}

func (r *FakeResourceChange) GetComputed() map[string]interface{} {
	return r.ComputedReturns
}
