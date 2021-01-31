package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/drlau/akashi/internal/compare"
	comparecmd "github.com/drlau/akashi/pkg/cmd/compare"
	diffcmd "github.com/drlau/akashi/pkg/cmd/diff"
	versioncmd "github.com/drlau/akashi/pkg/cmd/version"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/utils"
)

const (
	createKey  = "create"
	destroyKey = "destroy"
	updateKey  = "update"
)

var (
	version = "dev"
)

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
		Use:          "akashi <command> <path to ruleset>",
		Short:        "Akashi / 証",
		Long:         `Validate "terraform plan" changes against a customizable ruleset`,
		Args:         cobra.ExactArgs(1),
		RunE:         run,
		SilenceUsage: true,
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

	cmd.AddCommand(comparecmd.NewCmdCompare())
	cmd.AddCommand(diffcmd.NewCmdDiff())
	cmd.AddCommand(versioncmd.NewCmdVersion(os.Stdout, version))

	return cmd
}

func run(_ *cobra.Command, args []string) error {
	comparers, err := compare.Comparers(args[0])
	if err != nil {
		return err
	}

	in, err := plan.NewResourcePlans(file, json)
	if err != nil {
		return err
	}

	if quiet {
		fmt.Fprintln(os.Stderr, `[WARN] -q is deprecated. Please run "akashi compare" instead.`)
		os.Exit(runCompare(in, comparers))
	}
	out := utils.NewOutput(noColor)

	fmt.Fprintln(os.Stderr, `[WARN] no command is deprecated. Please run "akashi diff" instead.`)
	os.Exit(runDiff(out, in, comparers))
	return nil
}

func runCompare(rc []plan.ResourcePlan, comparers map[string]compare.Comparer) int {
	createComparer, hasCreate := comparers[createKey]
	destroyComparer, hasDestroy := comparers[destroyKey]
	updateComparer, hasUpdate := comparers[updateKey]

	for _, r := range rc {
		if r.IsCreate() && hasCreate {
			if !createComparer.Compare(r) {
				return 1
			}
		} else if r.IsDelete() && hasDestroy {
			if !destroyComparer.Compare(r) {
				return 1
			}
		} else if r.IsUpdate() && hasUpdate {
			if !updateComparer.Compare(r) {
				return 1
			}
		} else if strict {
			return 1
		}
	}

	return 0
}

func runDiff(out io.Writer, rc []plan.ResourcePlan, comparers map[string]compare.Comparer) int {
	exitCode := 0
	createComparer, hasCreate := comparers[createKey]
	destroyComparer, hasDestroy := comparers[destroyKey]
	updateComparer, hasUpdate := comparers[updateKey]

	for _, r := range rc {
		diff := ""
		pass := true
		if r.IsCreate() && hasCreate {
			diff, pass = createComparer.Diff(r)
		} else if r.IsDelete() && hasDestroy {
			diff, pass = destroyComparer.Diff(r)
		} else if r.IsUpdate() && hasUpdate {
			diff, pass = updateComparer.Diff(r)
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
