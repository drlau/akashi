package plan

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
