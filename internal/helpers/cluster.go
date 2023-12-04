package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ocClusterUrls struct {
	Server   string `json:"server"`
	ProxyUrl string `json:"proxy-url"`
}

type ocCluster struct {
	Name        string        `json:"name"`
	ClusterUrls ocClusterUrls `json:"cluster"`
}

type ocContext struct {
	Name    string            `json:"name"`
	Context map[string]string `json:"context"`
}

type OcConfig struct {
	Contexts       []ocContext `json:"contexts"`
	CurrentContext string      `json:"current-context"`
	Clusters       []ocCluster `json:"clusters"`
}

// Gets the current OpenShift cluster that a user is logged in.
func OcGetCurrentOcmCluster() (string, error) {
	ocmCluster := strings.TrimSpace(os.Getenv("OCM_CLUSTER"))
	return ocmCluster, nil
}

// Gets the current OpenShift namespace.
func OcGetCurrentNamespace(runAsOcUser string) (string, error) {
	config, err := OcGetConfig(runAsOcUser)
	if err != nil {
		return "", err
	}

	currentContext := config.CurrentContext

	for _, context := range config.Contexts {
		if context.Name == currentContext {
			return context.Context["namespace"], nil
		}
	}

	err = fmt.Errorf("current context not found: %s", currentContext)
	return "", err
}

func OcGetConfig(runAsOcUser string) (*OcConfig, error) {
	commandName := "oc"
	var commandArgs []string

	if len(runAsOcUser) > 0 {
		commandName = "sudo"
		commandArgs = []string{"-Eu", runAsOcUser, "oc", "config", "view", "-o", "json"}
	} else {
		commandArgs = []string{"config", "view", "-o", "json"}
	}

	bytes, err := exec.Command(
		commandName,
		commandArgs...).Output()

	if err != nil {
		return nil, err
	}

	var config OcConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func OcmGetOCMToken(
	ocmEnv string,
	ocmCliProd string,
	ocmCliStage string,
) (string, error) {
	var ocmToken string
	var err error
	var out []byte
	if ocmEnv == "production" && len(ocmCliProd) > 0 {
		out, err = exec.Command("sh", ocmCliProd, "token").CombinedOutput()
	} else if ocmEnv == "staging" && len(ocmCliStage) > 0 {
		out, err = exec.Command("sh", ocmCliStage, "token").CombinedOutput()
	} else {
		out, err = exec.Command("ocm", "token").CombinedOutput()
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", out, err)
	}

	ocmToken = string(out)
	return ocmToken, err
}
