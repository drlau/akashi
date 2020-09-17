package plan

import (
	"io"
	"io/ioutil"

	"github.com/hashicorp/terraform-json"
)

type jsonPlanChange struct {
	ResourceChange *tfjson.ResourceChange
}

func NewResourcePlanFromJSON(in io.Reader) ([]ResourceChange, error) {
	var result []ResourceChange

	parsed := &tfjson.Plan{}
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return result, err
	}

	err = parsed.UnmarshalJSON(data)
	if err != nil {
		return result, err
	}

	for _, rc := range parsed.ResourceChanges {
		result = append(result, newJSONPlanChange(rc))
	}

	return result, nil
}

func newJSONPlanChange(json *tfjson.ResourceChange) ResourceChange {
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
