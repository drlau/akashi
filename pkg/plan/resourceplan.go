package plan

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/drlau/tfplanparse"
	"github.com/hashicorp/terraform-json"
)

type ResourcePlan interface {
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

func NewResourcePlans(path string, isJSON bool) ([]ResourcePlan, error) {
	var result []ResourcePlan
	var data io.Reader
	var err error

	if path != "" {
		data, err = os.Open(path)
		if err != nil {
			return result, err
		}
	} else {
		data = os.Stdin
	}

	if isJSON {
		return NewResourcePlansFromJSON(data)
	}

	return NewResourcePlansFromPlanOutput(data)
}

func NewResourcePlansFromPlanOutput(in io.Reader) ([]ResourcePlan, error) {
	var result []ResourcePlan

	parsed, err := tfplanparse.Parse(in)
	if err != nil {
		return result, err
	}

	for _, rc := range parsed {
		result = append(result, NewTFPlanChange(rc))
	}

	return result, nil
}

func NewResourcePlansFromJSON(in io.Reader) ([]ResourcePlan, error) {
	var result []ResourcePlan

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
		result = append(result, NewJSONPlanChange(rc))
	}

	return result, nil
}
