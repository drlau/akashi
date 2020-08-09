package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/drlau/akashi/pkg/compare"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/ruleset"
)

const (
	createKey  = "create"
	destroyKey = "destroy"
)

var (
	file  string
	quiet bool
	json  bool
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "akashi <path to ruleset>",
		Short: "Akashi / 証",
		Long:  "Validate your terraform plan changes against a customizable ruleset",
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "read terraform json from file")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "don't output a diff")
	cmd.Flags().BoolVarP(&json, "json", "j", false, "read the contents as the output from 'terraform state show -json'")

	return cmd
}

func run(_ *cobra.Command, args []string) error {
	rulesetFile, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	var rs ruleset.Ruleset
	err = yaml.Unmarshal(rulesetFile, &rs)
	if err != nil {
		return err
	}

	var in []plan.ResourceChange
	var data io.Reader

	if file != "" {
		data, err = os.Open(file)
		if err != nil {
			return err
		}
	} else {
		data = os.Stdin
	}

	if json {
		in, err = plan.NewResourcePlanFromJSON(data)
		if err != nil {
			return err
		}
	} else {
		in, err = plan.NewResourcePlanFromPlanOutput(data)
		if err != nil {
			return err
		}
	}

	comparers := make(map[string]compare.Comparer)
	if rs.CreatedResources != nil {
		comparers[createKey] = compare.NewCreateComparer(*rs.CreatedResources)
	}
	if rs.DestroyedResources != nil {
		comparers[destroyKey] = compare.NewDestroyComparer(*rs.DestroyedResources)
	}

	if quiet {
		// TODO: better way to do this?
		os.Exit(runCompare(in, comparers, false))
	}

	runDiff(in, comparers, false)
	return nil
}

func runCompare(rc []plan.ResourceChange, comparers map[string]compare.Comparer, strict bool) int {
	createComparer, hasCreate := comparers[createKey]
	destroyComparer, hasDestroy := comparers[destroyKey]

	for _, r := range rc {
		if r.IsCreate() && hasCreate {
			if !createComparer.Compare(r) {
				return 1
			}
		} else if r.IsDelete() && hasDestroy {
			if !destroyComparer.Compare(r) {
				return 1
			}
		} else if strict {
			return 1
		}
	}

	return 0
}

func runDiff(rc []plan.ResourceChange, comparers map[string]compare.Comparer, strict bool) {
	createComparer, hasCreate := comparers[createKey]
	destroyComparer, hasDestroy := comparers[destroyKey]

	// TODO: handle output better
	var buf bytes.Buffer
	for _, r := range rc {
		if r.IsCreate() && hasCreate {
			result := createComparer.Diff(r)
			if result != "" {
				buf.WriteString(result)
			}
		} else if r.IsDelete() && hasDestroy {
			result := destroyComparer.Diff(r)
			if result != "" {
				buf.WriteString(result)
			}
		} else if strict {
			buf.WriteString(fmt.Sprintf("no comparer for %s\n", r.GetAddress()))
		}
	}
	fmt.Println(buf.String())
}
