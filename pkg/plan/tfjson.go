package plan

import (
	"github.com/hashicorp/terraform-json"
)

type jsonPlanChange struct {
	ResourceChange *tfjson.ResourceChange
}

func NewJSONPlanChange(json *tfjson.ResourceChange) *jsonPlanChange {
	return &jsonPlanChange{
		ResourceChange: json,
	}
}

func (j *jsonPlanChange) IsCreate() bool {
	return j.ResourceChange.Change.Actions.Create()
}

func (j *jsonPlanChange) IsDelete() bool {
	return j.ResourceChange.Change.Actions.Delete()
}

func (j *jsonPlanChange) IsNoOp() bool {
	return j.ResourceChange.Change.Actions.NoOp()
}

func (j *jsonPlanChange) IsUpdate() bool {
	return j.ResourceChange.Change.Actions.Update()
}

func (j *jsonPlanChange) GetBefore() map[string]interface{} {
	if j.ResourceChange.Change.Before != nil {
		return j.ResourceChange.Change.Before.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonPlanChange) GetAfter() map[string]interface{} {
	if j.ResourceChange.Change.After != nil {
		return j.ResourceChange.Change.After.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonPlanChange) GetBeforeChangedOnly() map[string]interface{} {
	return j.GetBefore()
}

func (j *jsonPlanChange) GetAfterChangedOnly() map[string]interface{} {
	return j.GetAfter()
}

func (j *jsonPlanChange) GetComputed() map[string]interface{} {
	if j.ResourceChange.Change.AfterUnknown != nil {
		return j.ResourceChange.Change.AfterUnknown.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (j *jsonPlanChange) GetName() string {
	return j.ResourceChange.Name
}

func (j *jsonPlanChange) GetType() string {
	return j.ResourceChange.Type
}

func (j *jsonPlanChange) GetAddress() string {
	return j.ResourceChange.Address
}
