package version

import (
	"fmt"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/spf13/cobra"
)

func NewVersionCommand(parentCmd *cobra.Command) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:           "version",
		Short:         "Show escope version information",
		Long:          "Display the current version, build date, and git commit hash of escope",
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("escope version %s\n", constants.Version)
		},
	}

	parentCmd.AddCommand(versionCmd)
	return versionCmd
}
