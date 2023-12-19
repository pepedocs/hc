package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	pkgInt "hc/internal"
	pkgIntHelper "hc/internal/helpers"
)

var (
	consoleCmdArgs struct {
		consoleContainerName string
		consoleContainerPort string
		ocmEnvironment       string
	}
)

var consoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Launches an OpenShift console for an existing hc container logged in to an OpenShift cluster.",
	Run:    launchOpenShiftConsole,
	PreRun: pkgInt.ToggleDebug,
}

type ocDeploymentContainer struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ocDeploymentTemplateSpec struct {
	Containers []ocDeploymentContainer `json:"containers"`
}

type ocDeploymentTemplate struct {
	Spec ocDeploymentTemplateSpec `json:"spec"`
}

type ocDeploymentSpec struct {
	Template ocDeploymentTemplate `json:"template"`
}
type ocDeployment struct {
	Spec ocDeploymentSpec `json:"spec"`
}

func launchOpenShiftConsole(cmd *cobra.Command, args []string) {
	config := pkgInt.GetHcConfig()
	ocUser := config.OcUser
	userHome := config.UserHome
	ocmCliAlias := config.OcmCliAlias

	ocmEnvironment := "production"
	if len(consoleCmdArgs.ocmEnvironment) > 0 {
		ocmEnvironment = loginCmdArgs.ocmEnvironment
	}

	var err error
	var out []byte
	switch {
	case ocmEnvironment == "production" && len(ocmCliAlias.OcmProduction) > 0:
		out, err = pkgIntHelper.RunCommandPipeStdin("sh", ocmCliAlias.OcmProduction, "post", "/api/accounts_mgmt/v1/access_token")
	case ocmEnvironment == "staging" && len(ocmCliAlias.OcmStaging) > 0:
		out, err = pkgIntHelper.RunCommandPipeStdin("sh", ocmCliAlias.OcmStaging, "post", "/api/accounts_mgmt/v1/access_token")
	default:
		out, err = pkgIntHelper.RunCommandPipeStdin("ocm", "post", "/api/accounts_mgmt/v1/access_token")

	}
	if err != nil {
		logger.Fatal("Failed to get access token: ", err)
	}

	path := fmt.Sprintf("%s/.kube/ocm-pull-secret-config.json", userHome)
	logger.Debugf("ocm-pull-secret path: %s", path)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Fatal("Failed to open ocm-pull-secret: ", err)
	}
	defer file.Close()
	file.WriteString(string(out))

	containerName := fmt.Sprintf("%s-openshift-console", consoleCmdArgs.consoleContainerName)
	kubeConfigFileName := fmt.Sprintf("%s/.kube/ocm-pull-secret-config.json", userHome)
	consoleListenAddr := fmt.Sprintf("http://0.0.0.0:%s", consoleCmdArgs.consoleContainerPort)
	pullArgs := []string{"pull", "--quiet", "--authfile", kubeConfigFileName}
	runArgs := []string{
		"run",
		"--rm",
		"--network",
		fmt.Sprintf("container:%s", consoleCmdArgs.consoleContainerName),
		"-e",
		"HTTPS_PROXY=http://squid.corp.redhat.com:3128",
		"--name",
		containerName,
		"--authfile",
		kubeConfigFileName,
	}

	out, err = pkgIntHelper.RunCommandOutput(
		"podman",
		"exec",
		"-it",
		"--user",
		ocUser,
		consoleCmdArgs.consoleContainerName,
		"oc",
		"get",
		"deployment",
		"console",
		"-n",
		"openshift-console",
		"-o",
		"json",
	)
	if err != nil {
		logger.Fatal("Failed to run command: ", err)
	}
	var openShiftConsoleDeploy ocDeployment
	var consoleImage string
	err = json.Unmarshal(out, &openShiftConsoleDeploy)
	if err != nil {
		logger.Fatal("Failed to unmarshal: ", err)
	}

	for _, container := range openShiftConsoleDeploy.Spec.Template.Spec.Containers {
		if container.Name == "console" {
			consoleImage = container.Image
			break
		}
	}

	out, err = pkgIntHelper.RunCommandOutput(
		"podman",
		"exec",
		"-it",
		"--user",
		ocUser,
		consoleCmdArgs.consoleContainerName,
		"oc",
		"config",
		"view",
		"-o",
		"json",
	)
	if err != nil {
		logger.Fatal("Failed to run command: ", err)
	}

	var clusterConfig pkgIntHelper.OcConfig
	err = json.Unmarshal(out, &clusterConfig)
	if err != nil {
		logger.Fatal("Failed to unmarshal: ", err)
	}

	imagePullArgs := append(pullArgs, consoleImage)
	_, err = pkgIntHelper.RunCommandOutput(
		"podman",
		imagePullArgs...,
	)
	if err != nil {
		logger.Fatal("Failed to run command: ", err)
	}
	cluster := clusterConfig.Clusters[0]
	apiUrl := cluster.ClusterUrls.Server
	alertManagerUrl := strings.Replace(apiUrl, "/backplane/cluster", "/backplane/alertmanager", 1)
	thanosUrl := strings.Replace(apiUrl, "/backplane/cluster", "/backplane/thanos", 1)
	alertManagerUrl = strings.TrimRight(alertManagerUrl, "/")
	thanosUrl = strings.TrimRight(thanosUrl, "/")

	out, err = pkgIntHelper.RunCommandOutput(
		"podman",
		"exec",
		"-it",
		"--user",
		ocUser,
		consoleCmdArgs.consoleContainerName,
		"ocm",
		"token",
	)
	if err != nil {
		logger.Fatal("Failed to run command: ", err)
	}
	ocmToken := strings.TrimSpace(string(out))
	baseAddress := fmt.Sprintf("http://127.0.0.1:%s", consoleCmdArgs.consoleContainerPort)
	runArgs = append(
		runArgs,
		consoleImage,
		"/opt/bridge/bin/bridge",
		"--public-dir",
		"/opt/bridge/static",
		"-base-address",
		baseAddress,
		"-branding",
		"dedicated",
		"-documentation-base-url",
		"https://docs.openshift.com/dedicated/4/",
		"-user-settings-location",
		"localstorage",
		"-user-auth",
		"disabled",
		"-k8s-mode",
		"off-cluster",
		"-k8s-auth",
		"bearer-token",
		"-k8s-mode-off-cluster-endpoint",
		apiUrl,
		"-k8s-mode-off-cluster-alertmanager",
		alertManagerUrl,
		"-k8s-mode-off-cluster-thanos",
		thanosUrl,
		"-k8s-auth-bearer-token",
		ocmToken,
		"-listen",
		consoleListenAddr,
		"-v",
		"5",
	)
	pkgIntHelper.RunCommandWithOsFiles("podman", os.Stdout, os.Stderr, os.Stdin, runArgs...)

}

func init() {
	rootCmd.AddCommand(consoleCmd)

	flags := consoleCmd.Flags()
	flags.StringVarP(
		&consoleCmdArgs.consoleContainerName,
		"consoleContainerName",
		"c",
		"",
		"The hc container name that is logged into an OpenShift cluster.",
	)

	flags.StringVarP(
		&consoleCmdArgs.consoleContainerPort,
		"consoleContainerPort",
		"p",
		"",
		"The hc container port that is logged into an OpenShift cluster.",
	)

	flags.StringVarP(
		&loginCmdArgs.ocmEnvironment,
		"ocmEnvironment",
		"e",
		"production",
		"OCM environemnt (production, staging)",
	)

	consoleCmd.MarkFlagRequired("consoleContainerName")
	consoleCmd.MarkFlagRequired("consoleContainerPort")
}
