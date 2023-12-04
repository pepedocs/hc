package cmd

import (
	"fmt"

	pkgIntHelper "hc/internal/helpers"

	"github.com/spf13/cobra"
)

// currentClusterCmd represents the currentCluster command
var currentClusterCmd = &cobra.Command{
	Use:   "currentCluster",
	Short: "Shows the current cluster where a user is logged in.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkContainerCommand(); err != nil {
			return
		}

		cluster, err := pkgIntHelper.OcGetCurrentOcmCluster()
		if err != nil {
			fmt.Print("")
		}
		fmt.Print(cluster)
	},
}

func init() {
	rootCmd.AddCommand(currentClusterCmd)
}
