package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func NewCmdVersion(out io.Writer, version string) *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, version)
		},
	}
}
