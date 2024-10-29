package validate

import (
	"fmt"
	"os"

	"github.com/drlau/akashi/internal/validate"
	"github.com/drlau/akashi/pkg/ruleset"

	"github.com/spf13/cobra"
)

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
			if err != nil {
				return fmt.Errorf("Could not parse ruleset: %v", err)
			}
			if res := validate.Validate(ruleset); !res.IsValid() {
				return fmt.Errorf("%s", res.String())
			}
			fmt.Println("Ruleset is valid!")
			return nil
		},
	}
	return cmd
}
