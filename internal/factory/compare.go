package factory

import (
	"io/ioutil"

	"github.com/drlau/akashi/pkg/compare"
	"github.com/drlau/akashi/pkg/ruleset"
	yaml "gopkg.in/yaml.v2"
)

const (
	CreateKey  = "create"
	DestroyKey = "destroy"
	UpdateKey  = "update"
)

// TODO: re-arrange to make interface here

func Comparers(path string) (map[string]compare.Comparer, error) {
	comparers := make(map[string]compare.Comparer)
	rulesetFile, err := ioutil.ReadFile(path)
	if err != nil {
		return comparers, err
	}

	var rs ruleset.Ruleset
	err = yaml.Unmarshal(rulesetFile, &rs)
	if err != nil {
		return comparers, err
	}

	if rs.CreatedResources != nil {
		comparers[CreateKey] = compare.NewCreateComparer(*rs.CreatedResources)
	}
	if rs.DestroyedResources != nil {
		comparers[DestroyKey] = compare.NewDestroyComparer(*rs.DestroyedResources)
	}
	if rs.UpdatedResources != nil {
		comparers[UpdateKey] = compare.NewUpdateComparer(*rs.UpdatedResources)
	}

	return comparers, nil
}
