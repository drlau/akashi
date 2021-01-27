package factory

import (
	"io"
	"os"

	"github.com/drlau/akashi/pkg/plan"
)

// TODO: re-arrange to make interface here

func ResourcePlans(path string, isJSON bool) ([]plan.ResourceChange, error) {
	var result []plan.ResourceChange
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
		return plan.NewResourcePlanFromJSON(data)
	}

	return plan.NewResourcePlanFromPlanOutput(data)
}
