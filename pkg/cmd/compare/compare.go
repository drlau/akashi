package compare

import (
	"fmt"

	"github.com/drlau/akashi/internal/compare"
	"github.com/drlau/akashi/pkg/plan"

	"github.com/spf13/cobra"
)

type CompareOptions struct {
	File   string
	JSON   bool
	Strict bool
}

func NewCmdCompare() *cobra.Command {
	opts := &CompareOptions{}
	cmd := &cobra.Command{
		Use:   "compare <path to ruleset>",
		Short: "Validate silently",
		Long:  `Validate "terraform plan" changes against a ruleset, exiting with code 0 if ok`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			comparers, err := compare.NewComparerSet(args[0])
			if err != nil {
				return err
			}

			plans, err := plan.NewResourcePlans(opts.File, opts.JSON)
			if err != nil {
				return err
			}

			cmd.SilenceErrors = true
			if result := runCompare(plans, comparers, opts.Strict); result != 0 {
				return fmt.Errorf("compare failed")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "read plan output from file")
	cmd.Flags().BoolVarP(&opts.Strict, "strict", "s", false, "require all resources to match a comparer")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "read the contents as the output from 'terraform state show -json'")

	return cmd
}

func runCompare(rc []plan.ResourcePlan, comparers compare.ComparerSet, strict bool) int {
	createComparer := comparers.CreateComparer
	destroyComparer := comparers.DestroyComparer
	updateComparer := comparers.UpdateComparer

	for _, r := range rc {
		if r.IsCreate() && createComparer != nil {
			if !createComparer.Compare(r) {
				return 1
			}
		} else if r.IsDelete() && destroyComparer != nil {
			if !destroyComparer.Compare(r) {
				return 1
			}
		} else if r.IsUpdate() && updateComparer != nil {
			if !updateComparer.Compare(r) {
				return 1
			}
		} else if strict {
			return 1
		}
	}

	return 0
}
