package resource

type ResourceValues struct {
	Values   map[string]interface{}
	Computed map[string]interface{}
}

func (rv ResourceValues) GetCombined() map[string]interface{} {
	combined := rv.Values
	if combined == nil {
		combined = make(map[string]interface{})
	}

	for k, v := range rv.Computed {
		combined[k] = v
	}

	return combined
}
