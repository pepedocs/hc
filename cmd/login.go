package cmd

import (
	"fmt"
	pkgInt "hc/internal"
	pkgIntHelper "hc/internal/helpers"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	loginCmdArgs struct {
		cluster        string
		ocmEnvironment string
		isOcmLoginOnly bool
	}
)

var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "Runs the hc container and logs into OCM and optionally to a cluster",
	PreRun: pkgInt.ToggleDebug,
	Run:    login,
}

func login(cmd *cobra.Command, args []string) {
	ocmEnvironment := "production"

	if len(loginCmdArgs.ocmEnvironment) > 0 {
		ocmEnvironment = loginCmdArgs.ocmEnvironment
	}

	ceFactory := pkgInt.NewCeFactory(map[string]interface{}{
		"ceName": "podman",
	})

	ocmCluster := loginCmdArgs.cluster
	isOcmLoginOnly := loginCmdArgs.isOcmLoginOnly

	ce, err := ceFactory.Create()
	if err != nil {
		log.Fatal("Failed to create container engine: ", err)
	}

	config := pkgInt.GetHcConfig()
	ocmLongLivedTokenPath := config.OcmLongLivedTokenPath
	var ocmToken string

	if len(ocmLongLivedTokenPath) > 0 {
		content, err := os.ReadFile(ocmLongLivedTokenPath)
		if err != nil {
			log.Fatalf("Failed to open long lived token file: %v", content)
		}
		ocmToken = string(content)
	} else {
		ocmCliAlias := config.OcmCliAlias
		ocmToken, err = pkgIntHelper.OcmGetOCMToken(
			loginCmdArgs.ocmEnvironment,
			ocmCliAlias.OcmProduction,
			ocmCliAlias.OcmStaging,
		)
		if err != nil {
			log.Fatalf("%s:\n%v", pkgInt.ErrOCMTokenFetchMsg, err)
		}
	}

	// Path where backplane config is mounted in the container
	containerBackplaneConfigPath := "/backplane-config.json"
	// Path where hc config is mounted in the container
	hcConfigPath := "/.hc.yaml"

	// Allocate free port and map host port for OpenShift console
	ports, err := pkgIntHelper.GetFreePorts(1)
	if err != nil {
		log.Fatal("Failed to generate port for Openshift console: ", err)
	}
	openshiftConsolePort := strconv.Itoa(ports[0])

	// Gather values for the container's environment variables
	ce.AppendEnvVar("HOST_USER", config.HostUser)
	ce.AppendEnvVar("OC_USER", config.OcUser)
	ce.AppendEnvVar("OCM_CLUSTER", ocmCluster)
	ce.AppendEnvVar("IS_OCM_LOGIN_ONLY", strconv.FormatBool(isOcmLoginOnly))
	ce.AppendEnvVar("OCM_TOKEN", ocmToken)
	ce.AppendEnvVar("IS_IN_CONTAINER", "true")
	ce.AppendEnvVar("OCM_ENVIRONMENT", ocmEnvironment)
	ce.AppendEnvVar("BACKPLANE_CONFIG", containerBackplaneConfigPath)
	ce.AppendEnvVar("OPENSHIFT_CONSOLE_PORT", openshiftConsolePort)

	// Gather values for the container's host-mounted volumes
	if ocmEnvironment == "production" {
		ce.AppendVolMap(
			fmt.Sprintf("%s/.config/backplane/%s", config.UserHome, config.BackplaneConfigProd),
			containerBackplaneConfigPath,
			"ro",
		)
	} else {
		ce.AppendVolMap(
			fmt.Sprintf("%s/.config/backplane/%s", config.UserHome, config.BackplaneConfigStage),
			containerBackplaneConfigPath,
			"ro",
		)
	}
	ce.AppendVolMap(fmt.Sprintf("%s/.hc.yaml", config.UserHome), hcConfigPath, "ro")

	for _, dirMap := range config.CustomDirMaps {
		ce.AppendVolMap(dirMap.HostDir, dirMap.ContainerDir, dirMap.FileAttrs)
	}

	// Gather values for the containers host-mapped TCP ports
	// Openshift console port
	ce.AppendPortMap(openshiftConsolePort, openshiftConsolePort, "127.0.0.1")

	suffix := uuid.New()
	containerName := fmt.Sprintf("hc-%s-%s", ocmCluster, suffix.String()[:6])

	runCmd := ce.GetRunCmd(
		containerName,
		"./hc",
		"hc:latest",
		"clusterLogin",
		ocmCluster,
		"--config",
		hcConfigPath,
	)

	if pkgInt.Debug {
		runCmd = append(runCmd, "-d")
	}

	log.Debugf("Container run command: %v", runCmd)

	pkgIntHelper.RunCommandWithOsFiles(
		ce.GetExecName(),
		os.Stdout,
		os.Stderr,
		os.Stdin,
		runCmd...,
	)

}

func init() {
	rootCmd.AddCommand(loginCmd)

	flags := loginCmd.Flags()
	flags.StringVarP(
		&loginCmdArgs.cluster,
		"ocmCluster",
		"c",
		"",
		"Cluster name or id.",
	)

	flags.StringVarP(
		&loginCmdArgs.ocmEnvironment,
		"ocmEnvironment",
		"e",
		"production",
		"OCM environemnt (production, staging)",
	)

	flags.BoolVar(
		&loginCmdArgs.isOcmLoginOnly,
		"isOcmLoginOnly",
		false,
		"Log in to OCM only.",
	)
}
