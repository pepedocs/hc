package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	pkgIntHelper "hc/internal/helpers"
)

var (
	currentNamespaceCmdArgs struct {
		ocUser string
	}
)

var currentNamespaceCmd = &cobra.Command{
	Use:   "currentNamespace",
	Short: "Shows OpenShift's current context namespace given an OpenShift user.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkContainerCommand(); err != nil {
			return
		}

		namespace, err := pkgIntHelper.OcGetCurrentNamespace(currentNamespaceCmdArgs.ocUser)
		if err != nil {
			fmt.Print("")
		}
		fmt.Print(namespace)
	},
}

func init() {
	rootCmd.AddCommand(currentNamespaceCmd)
	currentNamespaceCmd.Flags().StringVarP(
		&currentNamespaceCmdArgs.ocUser,
		"ocUser",
		"u",
		"",
		"Run as OpenShift user.",
	)
	currentNamespaceCmd.MarkFlagRequired("ocUser")
}
