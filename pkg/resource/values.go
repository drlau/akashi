package resource

type ResourceValues struct {
	Values map[string]interface{}
	// TODO: better implementation of ChangedValues(a filter operation on Values seems ideal)
	ChangedValues map[string]interface{}
	Computed      map[string]interface{}
}

// TODO: fix merging maps
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
