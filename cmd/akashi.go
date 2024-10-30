package cmd

import (
	"os"

	"github.com/spf13/cobra"

	comparecmd "github.com/drlau/akashi/pkg/cmd/compare"
	diffcmd "github.com/drlau/akashi/pkg/cmd/diff"
	matchcmd "github.com/drlau/akashi/pkg/cmd/match"
	validatecmd "github.com/drlau/akashi/pkg/cmd/validate"
	versioncmd "github.com/drlau/akashi/pkg/cmd/version"
)

var version = "dev"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "akashi <command> <path to ruleset>",
		Short:        "Akashi / 証",
		Long:         `Validate "terraform plan" changes against a customizable ruleset`,
		SilenceUsage: true,
	}

	cmd.SetVersionTemplate(version)

	cmd.AddCommand(comparecmd.NewCmdCompare())
	cmd.AddCommand(diffcmd.NewCmdDiff())
	cmd.AddCommand(matchcmd.NewCmdMatch())
	cmd.AddCommand(validatecmd.NewCmd())
	cmd.AddCommand(versioncmd.NewCmdVersion(os.Stdout, version))

	return cmd
}
