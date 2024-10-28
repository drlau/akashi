package match

import (
	"fmt"
	"io"
	"strings"

	"github.com/drlau/akashi/internal/compare"
	"github.com/drlau/akashi/pkg/plan"
	"github.com/drlau/akashi/pkg/utils"

	"github.com/spf13/cobra"
)

type MatchOptions struct {
	File      string
	JSON      bool
	Invert    bool
	Separator string
}

func NewCmdMatch() *cobra.Command {
	opts := &MatchOptions{}
	cmd := &cobra.Command{
		Use:   "match <path to ruleset>",
		Short: "Outputs resource paths which match the ruleset",
		Long:  `Outputs resource paths from "terraform plan" which are defined in the ruleset`,
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

			out := utils.NewOutput(true)
			cmd.SilenceErrors = true
			runMatch(out, plan, comparers, opts)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "read plan output from file")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "read the contents as the output from 'terraform state show -json'")
	cmd.Flags().BoolVarP(&opts.Invert, "invert", "i", false, "outputs resources which do not match the ruleset")
	cmd.Flags().StringVarP(&opts.Separator, "separator", "s", "\n", "separator between resource paths")

	return cmd
}

func runMatch(out io.Writer, rc []plan.ResourcePlan, comparers compare.ComparerSet, opts *MatchOptions) {
	createComparer := comparers.CreateComparer
	destroyComparer := comparers.DestroyComparer
	updateComparer := comparers.UpdateComparer

	var matches []string
	for _, r := range rc {
		var match bool
		if r.IsCreate() && createComparer != nil {
			match = createComparer.Compare(r)
		} else if r.IsDelete() && destroyComparer != nil {
			match = destroyComparer.Compare(r)
		} else if r.IsUpdate() && updateComparer != nil {
			match = updateComparer.Compare(r)
		}
		if (!opts.Invert && match) || (opts.Invert && !match) {
			matches = append(matches, r.GetAddress())
		}
	}

	fmt.Fprintln(out, strings.Join(matches, opts.Separator))
}
