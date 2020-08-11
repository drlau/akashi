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

type FailedArg struct {
	Expected interface{}
	Actual   interface{}
}
