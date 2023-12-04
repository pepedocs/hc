package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pkgInt "hc/internal"
	pkgIntHelper "hc/internal/helpers"

	"github.com/spf13/cobra"

	logger "github.com/sirupsen/logrus"
)

var hcCon *hcContainer

var clusterLoginCmd = &cobra.Command{
	Use:    "clusterLogin",
	Short:  "Logs in to an hybrid-cloud OpenShift cluster.",
	PreRun: pkgInt.ToggleDebug,
	Run:    clusterLogin,
}

func clusterLogin(cmd *cobra.Command, args []string) {
	hcCon = NewHcContainer(pkgInt.GetHcConfig())
	if err := checkContainerCommand(); err != nil {
		logger.Fatal(err)
	}

	configureOCMUser()
	configureWorkspaceDirs()
	OCMLogin()
	OCMBackplaneLogin()

	customPortMapsStr := strings.Trim(getEnvVar("CUSTOM_PORT_MAPS"), ",")
	var allocatedContainerPorts []string

	if len(customPortMapsStr) > 0 {
		for _, pm := range strings.Split(customPortMapsStr, ",") {
			ports := strings.Split(pm, ":")
			allocatedContainerPorts = append(allocatedContainerPorts, ports[1])
		}
	}

	runTerminal()
}

func configureOCMUser() {
	// Configure ocm user
	commands := [][]string{
		{
			"useradd",
			"-m",
			hcCon.HostUser,
			"-d",
			hcCon.UserHome,
		},
		{
			"usermod",
			"-aG",
			"wheel",
			hcCon.HostUser,
		},
	}
	errors := pkgIntHelper.RunCommandListStreamOutput(commands)

	if len(errors) > 0 {
		logger.Fatalf("Encountered errors while configuring OCM user: %v", errors)
	}

	line := []byte("\n%wheel         ALL = (ALL) NOPASSWD: ALL\n")
	os.WriteFile("/etc/sudoer", line, 0644)

}

func configureWorkspaceDirs() {
	// Configure directories
	commands := [][]string{
		{
			"mkdir",
			"-p",
			fmt.Sprintf("%s/.kube", hcCon.UserHome),
		},
		{
			"chown",
			"-R",
			fmt.Sprintf("%s:%s", hcCon.HostUser, hcCon.HostUser),
			fmt.Sprintf("%s/.kube", hcCon.UserHome),
		},
		{
			"mkdir",
			"-p",
			fmt.Sprintf("%s/.config/ocm", hcCon.UserHome),
		},
		{
			"chown",
			"-R",
			fmt.Sprintf("%s:%s", hcCon.HostUser, hcCon.HostUser),
			fmt.Sprintf("%s/.config/ocm", hcCon.UserHome),
		},
		{
			"chmod",
			"o+rwx",
			"/hc",
		},
	}
	errors := pkgIntHelper.RunCommandListStreamOutput(commands)

	if len(errors) > 0 {
		logger.Fatalf("Encountered errors while configuring hc directories: %v", errors)
	}

}

func OCMLogin() {
	logger.Info("Logging into ocm ", hcCon.OcmEnvironment)

	status := pkgIntHelper.RunCommandStreamOutput(
		"sudo",
		"-Eu",
		hcCon.HostUser,
		"ocm",
		"login",
		fmt.Sprintf("--token=%s", hcCon.OcmToken),
		fmt.Sprintf("--url=%s", hcCon.OcmEnvironment),
	)

	if status.Exit != 0 {
		logger.Fatalf("OCM Login failed: %v", status.Error)
	}

	logger.Info("OCM Login successful.")
}

func OCMBackplaneLogin() {
	isOcmLoginOnly, err := strconv.ParseBool(hcCon.IsOcmLoginOnly)
	if err != nil {
		logger.Fatal("Failed to parse environment variable: ", err)
	}

	if !isOcmLoginOnly {
		// Backplane login
		status := pkgIntHelper.RunCommandStreamOutput(
			"sudo",
			"-Eu",
			hcCon.HostUser,
			"ocm",
			"backplane",
			"login",
			hcCon.OcmCluster,
		)

		if status.Exit != 0 {
			logger.Fatalf("OCM backplane login failed: %v", status.Error)
		}
		logger.Info("OCM backplane login successful.")
	}

}

func runTerminal() {
	// Run terminal
	status := pkgIntHelper.RunCommandStreamOutput("cp", "/terminal/bashrc", hcCon.UserBashrcPath)
	if status.Exit != 0 {
		logger.Fatalf("Failed to copy /terminal/bashrc: %v", status.Error)
	}

	file, err := os.OpenFile(hcCon.UserBashrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf("Failed to open %s: %s\n", hcCon.UserBashrcPath, err)
	} else {
		defer file.Close()

		ps1String := fmt.Sprintf(
			"\nPS1='[%s %s $(/usr/bin/hc currentCluster) $(/usr/bin/hc currentNamespace -u %s)]$ '\n",
			hcCon.HostUser,
			hcCon.OcmEnvironment,
			hcCon.HostUser)
		_, err = file.WriteString(ps1String)
		if err != nil {
			logger.Errorf("Failed to write to file %s: %s\n", hcCon.UserBashrcPath, err)
		}

		config := pkgInt.GetHcConfig()
		exportStr := "\nexport PATH=$PATH"
		for _, path := range config.AddToPATHEnv {
			exportStr += fmt.Sprintf(":%s", path)
		}
		_, err = file.WriteString(exportStr)
		if err != nil {
			logger.Errorf("Failed to write to file %s: %s\n", hcCon.UserBashrcPath, err)
		}

		for _, path := range config.ExportEnvVars {
			exportStr = fmt.Sprintf("\nexport %s", path)
			_, err = file.WriteString(exportStr)
			if err != nil {
				logger.Errorf("Failed to write to file %s: %s\n", hcCon.UserBashrcPath, err)
			}
		}
	}
	pkgIntHelper.RunCommandWithOsFiles("sudo", os.Stdout, os.Stderr, os.Stdin, "-Eu", hcCon.HostUser, "bash")
}

func init() {
	rootCmd.AddCommand(clusterLoginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterLoginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterLoginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
