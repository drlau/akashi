package resource

type CompareResult struct {
	// args that had a matching EnforcedValue and are equal
	Enforced map[string]interface{}

	// args that had a matching EnforcedValue but were not equal
	Failed map[string]interface{}

	// args that had a matching Ignored entry
	Ignored map[string]interface{}

	// args that did not have an enforcedValue and were not ignored
	Extra map[string]interface{}

	// args defined in EnforcedValue but were not present
	MissingEnforced map[string]interface{}

	// args defined in Ignored but were not present
	MissingIgnored map[string]interface{}
}

func (cr *CompareResult) GetEnforced() map[string]interface{} {
	return cr.Enforced
}

func (cr *CompareResult) GetFailed() map[string]interface{} {
	return cr.Failed
}

func (cr *CompareResult) GetIgnored() map[string]interface{} {
	return cr.Ignored
}

func (cr *CompareResult) GetExtra() map[string]interface{} {
	return cr.Extra
}

func (cr *CompareResult) GetMissingEnforced() map[string]interface{} {
	return cr.MissingEnforced
}

func (cr *CompareResult) GetMissingIgnored() map[string]interface{} {
	return cr.MissingIgnored
}

type FailedArg struct {
	Expected interface{}
	Actual   interface{}
	MatchAny bool
}
