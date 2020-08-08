package plan

import (
	"github.com/hashicorp/terraform-json"
)

type jsonChange struct {
	ResourceChange *tfjson.ResourceChange
}

func NewResourceChangeFromJSON(json *tfjson.ResourceChange) ResourceChange {
	return &jsonChange{
		ResourceChange: json,
	}
}

func (j *jsonChange) IsCreate() bool {
	return j.ResourceChange.Change.Actions.Create()
}

func (j *jsonChange) IsDelete() bool {
	return j.ResourceChange.Change.Actions.Delete()
}

func (j *jsonChange) GetBefore() map[string]interface{} {
	if j.ResourceChange.Change.Before != nil {
		return j.ResourceChange.Change.Before.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonChange) GetAfter() map[string]interface{} {
	if j.ResourceChange.Change.After != nil {
		return j.ResourceChange.Change.After.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonChange) GetComputed() map[string]interface{} {
	if j.ResourceChange.Change.AfterUnknown != nil {
		return j.ResourceChange.Change.AfterUnknown.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonChange) GetName() string {
	return j.ResourceChange.Name
}

func (j *jsonChange) GetType() string {
	return j.ResourceChange.Type
}

func (j *jsonChange) GetAddress() string {
	return j.ResourceChange.Address
}
