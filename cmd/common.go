package cmd

import (
	"errors"
	"os"
	"strings"
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
