package fakes

type FakeResourcePlan struct {
	AddressReturns  string
	NameReturns     string
	TypeReturns     string
	CreateReturns   bool
	DeleteReturns   bool
	NoOpReturns     bool
	UpdateReturns   bool
	BeforeReturns   map[string]interface{}
	AfterReturns    map[string]interface{}
	ComputedReturns map[string]interface{}
}

func (r *FakeResourcePlan) GetAddress() string {
	return r.AddressReturns
}

func (r *FakeResourcePlan) GetName() string {
	return r.NameReturns
}

func (r *FakeResourcePlan) GetType() string {
	return r.TypeReturns
}

func (r *FakeResourcePlan) IsCreate() bool {
	return r.CreateReturns
}

func (r *FakeResourcePlan) IsDelete() bool {
	return r.DeleteReturns
}

func (r *FakeResourcePlan) IsNoOp() bool {
	return r.NoOpReturns
}

func (r *FakeResourcePlan) IsUpdate() bool {
	return r.UpdateReturns
}

func (r *FakeResourcePlan) GetBefore() map[string]interface{} {
	return r.BeforeReturns
}

func (r *FakeResourcePlan) GetAfter() map[string]interface{} {
	return r.AfterReturns
}

func (r *FakeResourcePlan) GetBeforeChangedOnly() map[string]interface{} {
	return r.BeforeReturns
}

func (r *FakeResourcePlan) GetAfterChangedOnly() map[string]interface{} {
	return r.AfterReturns
}

func (r *FakeResourcePlan) GetComputed() map[string]interface{} {
	return r.ComputedReturns
}
