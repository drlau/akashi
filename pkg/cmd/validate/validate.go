package validate

import (
	"os"
	"fmt"

	"github.com/drlau/akashi/internal/validate"
	"github.com/drlau/akashi/pkg/ruleset"

	"github.com/spf13/cobra"
)

// TODO: print out which resources are not valid
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <path to ruleset>",
		Short: "Validte the ruleset",
		Long:  "Validate the ruleset, exiting with code 0 if the ruleset is valid",

		// NOTE: We explicitly do not set Args with ExactArgs(1) since that
		// will not print the help message.
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				cmd.Help()
				os.Exit(1)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleset, err := ruleset.ParseRuleset(args[0])
			if res := validate.Validate(ruleset); !res.IsValid() {
				return fmt.Errorf("ruleset is invalid")
			}
			return err
		},
	}
	return cmd
}
