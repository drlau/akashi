package resource

type ResourceValues struct {
	Values   map[string]interface{}
	Computed map[string]interface{}
}

func (rv ResourceValues) GetCombined() map[string]interface{} {
	return setUnion(rv.Values, rv.Computed)
}
