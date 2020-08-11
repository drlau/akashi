package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/drlau/akashi/pkg/compare"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/ruleset"
	"github.com/drlau/akashi/pkg/utils"
)

const (
	createKey  = "create"
	destroyKey = "destroy"
)

// TODO: set this dynamically
const version = "0.0.2"

var (
	file          string
	versionOutput string
	quiet         bool
	json          bool
	failedOnly    bool
	strict        bool
	noColor       bool
	errorOnFail   bool
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "akashi <path to ruleset>",
		Short: "Akashi / 証",
		Long:  `Validate "terraform plan" changes against a customizable ruleset`,
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}

	cmd.SetVersionTemplate(version)

	cmd.Flags().StringVarP(&file, "file", "f", "", "read plan output from file")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "compare only, and error if there is a failing rule")
	cmd.Flags().BoolVar(&failedOnly, "failed-only", false, "only output failing lines")
	cmd.Flags().BoolVarP(&strict, "strict", "s", false, "require all resources to match a comparer")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "disable color output")
	cmd.Flags().BoolVarP(&errorOnFail, "error-on-fail", "e", false, "for non-quiet runs, make akashi return exit code 1 on fails")
	cmd.Flags().BoolVarP(&json, "json", "j", false, "read the contents as the output from 'terraform state show -json'")
	// TODO
	// cmd.Flags().BoolVarP(&verbose, "verbose", "V", false, "enable verbose output")

	versionCmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	cmd.AddCommand(versionCmd)

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
		os.Exit(runCompare(in, comparers))
	}
	out := utils.NewOutput(noColor)

	os.Exit(runDiff(out, in, comparers))
	return nil
}

func runCompare(rc []plan.ResourceChange, comparers map[string]compare.Comparer) int {
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

func runDiff(out io.Writer, rc []plan.ResourceChange, comparers map[string]compare.Comparer) int {
	exitCode := 0
	createComparer, hasCreate := comparers[createKey]
	destroyComparer, hasDestroy := comparers[destroyKey]

	for _, r := range rc {
		diff := ""
		pass := true
		if r.IsCreate() && hasCreate {
			diff, pass = createComparer.Diff(r)
		} else if r.IsDelete() && hasDestroy {
			diff, pass = destroyComparer.Diff(r)
		} else {
			if !strict {
				continue
			}

			if errorOnFail {
				exitCode = 1
			}
			fmt.Fprintln(out, fmt.Sprintf("%s %s (no matching comparer)", utils.Yellow("?"), r.GetAddress()))
			continue
		}
		if pass {
			if failedOnly {
				continue
			}

			fmt.Fprintln(out, diff)
			continue
		}

		fmt.Fprintln(out, diff)
		if errorOnFail {
			exitCode = 1
		}
	}

	return exitCode
}
