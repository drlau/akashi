package plan

import (
	"io"

	"github.com/drlau/tfplanparse"
)

type tfPlanChange struct {
	ResourceChange *tfplanparse.ResourceChange
}

func NewResourcePlanFromPlanOutput(in io.Reader) ([]ResourceChange, error) {
	var result []ResourceChange

	parsed, err := tfplanparse.Parse(in)
	if err != nil {
		return result, err
	}

	for _, rc := range parsed {
		result = append(result, newTFPlanChange(rc))
	}

	return result, nil
}

func newTFPlanChange(rc *tfplanparse.ResourceChange) ResourceChange {
	return &tfPlanChange{
		ResourceChange: rc,
	}
}

func (t *tfPlanChange) IsCreate() bool {
	return t.ResourceChange.UpdateType == tfplanparse.NewResource
}

func (t *tfPlanChange) IsDelete() bool {
	return t.ResourceChange.UpdateType == tfplanparse.DestroyResource
}

func (t *tfPlanChange) IsNoOp() bool {
	return t.ResourceChange.UpdateType == tfplanparse.NoOpResource
}

func (t *tfPlanChange) IsUpdate() bool {
	return t.ResourceChange.UpdateType == tfplanparse.UpdateInPlaceResource || t.ResourceChange.UpdateType == tfplanparse.ForceReplaceResource
}

func (t *tfPlanChange) GetBefore() map[string]interface{} {
	return t.ResourceChange.GetBeforeResource(tfplanparse.IgnoreSensitive)
}

func (t *tfPlanChange) GetAfter() map[string]interface{} {
	return t.ResourceChange.GetAfterResource(tfplanparse.IgnoreSensitive)
}

func (t *tfPlanChange) GetBeforeChangedOnly() map[string]interface{} {
	return t.ResourceChange.GetBeforeResource(tfplanparse.IgnoreSensitive, tfplanparse.IgnoreNoOp)
}

func (t *tfPlanChange) GetAfterChangedOnly() map[string]interface{} {
	return t.ResourceChange.GetAfterResource(tfplanparse.IgnoreSensitive, tfplanparse.IgnoreNoOp)
}

func (t *tfPlanChange) GetComputed() map[string]interface{} {
	return t.ResourceChange.GetAfterResource(tfplanparse.ComputedOnly)
}

func (t *tfPlanChange) GetName() string {
	return t.ResourceChange.Name
}

func (t *tfPlanChange) GetType() string {
	return t.ResourceChange.Type
}

func (t *tfPlanChange) GetAddress() string {
	return t.ResourceChange.Address
}
