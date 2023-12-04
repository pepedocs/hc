package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pkgInt "hc/internal"

	logger "github.com/sirupsen/logrus"
)

func isInContainer() bool {
	return os.Getenv("IS_IN_CONTAINER") == "true"
}

func checkContainerCommand() error {
	if !isInContainer() {
		return errors.New("this command is intended to be run only inside the workspace container")
	}
	return nil
}

func getEnvVar(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

type hcContainer struct {
	HostUser       string
	UserHome       string
	IsOcmLoginOnly string
	customPortMaps string
	UserBashrcPath string
	OcmCluster     string
	OcmToken       string
	OcmEnvironment string
}

func NewHcContainer(config *pkgInt.HcConfig) *hcContainer {
	if config == nil {
		logger.Fatal("Config is not yet available at this stage.")
	}
	return &hcContainer{
		HostUser:       getEnvVar("HOST_USER"),
		UserHome:       config.UserHome,
		IsOcmLoginOnly: getEnvVar("IS_OCM_LOGIN_ONLY"),
		customPortMaps: getEnvVar("CUSTOM_PORT_MAPS"),
		UserBashrcPath: fmt.Sprintf("%s/.bashrc", config.UserHome),
		OcmCluster:     getEnvVar("OCM_CLUSTER"),
		OcmToken:       getEnvVar("OCM_TOKEN"),
		OcmEnvironment: getEnvVar("OCM_ENVIRONMENT"),
	}
}
