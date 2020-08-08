package plan

type ResourceChange interface {
	IsCreate() bool
	IsDelete() bool
	GetBefore() map[string]interface{}
	GetAfter() map[string]interface{}
	GetComputed() map[string]interface{}
	GetName() string
	GetType() string
	GetAddress() string
}
