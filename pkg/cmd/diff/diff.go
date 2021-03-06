package diff

import (
	"fmt"
	"io"

	"github.com/drlau/akashi/internal/compare"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/utils"

	"github.com/spf13/cobra"
)

type DiffOptions struct {
	File        string
	JSON        bool
	FailedOnly  bool
	Strict      bool
	NoColor     bool
	ErrorOnFail bool
}

func NewCmdDiff() *cobra.Command {
	opts := &DiffOptions{}
	cmd := &cobra.Command{
		Use:   "diff <path to ruleset>",
		Short: "Validate changes",
		Long:  `Validate "terraform plan" changes against a ruleset`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			comparers, err := compare.NewComparerSet(args[0])
			if err != nil {
				return err
			}

			plan, err := plan.NewResourcePlans(opts.File, opts.JSON)
			if err != nil {
				return err
			}

			out := utils.NewOutput(opts.NoColor)
			cmd.SilenceErrors = true
			if result := runDiff(out, plan, comparers, opts); result != 0 {
				return fmt.Errorf("diff failed")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "read plan output from file")
	cmd.Flags().BoolVarP(&opts.Strict, "strict", "s", false, "require all resources to match a comparer")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "read the contents as the output from 'terraform state show -json'")
	cmd.Flags().BoolVar(&opts.FailedOnly, "failed-only", false, "only output failing lines")
	cmd.Flags().BoolVar(&opts.NoColor, "no-color", false, "disable color output")
	cmd.Flags().BoolVarP(&opts.ErrorOnFail, "error-on-fail", "e", false, "return exit code 1 on fail")

	return cmd
}

func runDiff(out io.Writer, rc []plan.ResourcePlan, comparers compare.ComparerSet, opts *DiffOptions) int {
	exitCode := 0
	createComparer := comparers.CreateComparer
	destroyComparer := comparers.DestroyComparer
	updateComparer := comparers.UpdateComparer

	for _, r := range rc {
		diff := ""
		pass := true
		if r.IsCreate() && createComparer != nil {
			diff, pass = createComparer.Diff(r)
		} else if r.IsDelete() && destroyComparer != nil {
			diff, pass = destroyComparer.Diff(r)
		} else if r.IsUpdate() && updateComparer != nil {
			diff, pass = updateComparer.Diff(r)
		} else {
			if !opts.Strict {
				continue
			}

			if opts.ErrorOnFail {
				exitCode = 1
			}
			fmt.Fprintln(out, fmt.Sprintf("%s %s (no matching comparer)", utils.Yellow("?"), r.GetAddress()))
			continue
		}
		if pass {
			if opts.FailedOnly {
				continue
			}

			fmt.Fprintln(out, diff)
			continue
		}

		fmt.Fprintln(out, diff)
		if opts.ErrorOnFail {
			exitCode = 1
		}
	}

	return exitCode
}
