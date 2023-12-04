package cmd

import (
	pkgInt "hc/internal"
	pkgIntHelper "hc/internal/helpers"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the hc image",
	Run: func(cmd *cobra.Command, args []string) {
		ceFactory := pkgInt.NewCeFactory(map[string]interface{}{
			"ceName": "podman",
		})

		ce, err := ceFactory.Create()
		if err != nil {
			log.Fatal("Failed to create container engine: ", err)
		}

		config := pkgInt.GetHcConfig()

		ce.AppendBuildArg("BASE_IMAGE_VERSION", config.BaseImageVersion)
		ce.AppendBuildArg("OCM_CLI_VERSION", config.OCMCLIVersion)
		ce.AppendBuildArg("BACKPLANE_CLI_VERSION", config.BackplaneCLIVersion)

		out, err := pkgIntHelper.RunCommandOutput(
			"git",
			"rev-parse",
			"HEAD",
		)
		if err != nil {
			log.Fatal(err)
		}
		headSha := string(out)
		ce.AppendBuildArg("BUILD_SHA", headSha)

		ceBuildCmd := ce.GetBuildCmd()
		pkgIntHelper.RunCommandStreamOutput(ce.GetExecName(), ceBuildCmd...)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
