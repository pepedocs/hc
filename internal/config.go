package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type DirMap struct {
	HostDir      string `mapstructure:"hostDir"`
	ContainerDir string `mapstructure:"containerDir"`
	FileAttrs    string `mapstructure:"fileAttrs"`
}

type PortMap struct {
	HostPort      string `mapstructure:"hostPort"`
	ContainerPort string `mapstructure:"containerPort"`
}

type OcmCliAlias struct {
	OcmProduction string `mapstructure:"production"`
	OcmStaging    string `mapstructure:"staging"`
}

type HcConfig struct {
	CustomDirMaps         []DirMap    `mapstructure:"customDirMaps"`
	AddToPATHEnv          []string    `mapstructure:"addToPATHEnv"`
	ExportEnvVars         []string    `mapstructure:"exportEnvVars"`
	HostUser              string      `mapstructure:"hostUser"`
	OcUser                string      `mapstructure:"ocUser"`
	UserHome              string      `mapstructure:"userHome"`
	BackplaneConfigProd   string      `mapstructure:"backplaneConfigProd"`
	BackplaneConfigStage  string      `mapstructure:"backplaneConfigStage"`
	BaseImageVersion      string      `mapstructure:"baseImageVersion"`
	OCMCLIVersion         string      `mapstructure:"ocmCLIVersion"`
	BackplaneCLIVersion   string      `mapstructure:"backplaneCLIVersion"`
	CustomPortMaps        []PortMap   `mapstructure:"customPortMaps"`
	OcmLongLivedTokenPath string      `mapstructure:"ocmLongLivedTokenPath"`
	OcmCliAlias           OcmCliAlias `mapstructure:"ocmCLIAlias"`
}

func GetHcConfig() *HcConfig {
	var conf HcConfig
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Failed to unmarshal config: ", err)
	}
	return &conf
}

func ValidateConfig() error {
	conf := GetHcConfig()
	requiredStrs := map[string]string{
		"hostUser":             conf.HostUser,
		"ocUser":               conf.OcUser,
		"userHome":             conf.UserHome,
		"backplaneConfigProd":  conf.BackplaneConfigProd,
		"backplaneConfigStage": conf.BackplaneConfigStage,
		"ocmCLIVersion":        conf.OCMCLIVersion,
		"backplaneCLIVersion":  conf.BackplaneCLIVersion,
	}

	for name, value := range requiredStrs {
		noSpaces := strings.TrimSpace(value)
		if len(noSpaces) == 0 {
			return fmt.Errorf("missing required config: \"%s\"", name)
		}
	}
	return nil
}

func (c *HcConfig) GetAddToPATHEnv() []string {
	return c.AddToPATHEnv
}

func (c *HcConfig) GetCustomDirMaps() []DirMap {
	return c.CustomDirMaps
}

func (c *HcConfig) GetExportEnvVars() []string {
	return c.ExportEnvVars
}

func (c *HcConfig) GetHostUser() string {
	return c.HostUser
}

func (c *HcConfig) GetOcUser() string {
	return c.OcUser
}

func (c *HcConfig) GetUserHome() string {
	return c.UserHome
}

func (c *HcConfig) GetBackplaneConfigProd() string {
	return c.BackplaneConfigProd
}

func (c *HcConfig) GetBackplaneConfigStage() string {
	return c.BackplaneConfigStage
}

func (c *HcConfig) GetBaseImage() string {
	return fmt.Sprintf("fedora:%s", c.BaseImageVersion)
}

func (c *HcConfig) GetOcmCLIVersion() string {
	return c.OCMCLIVersion
}

func (c *HcConfig) GetBackplaneCLIVersion() string {
	return c.BackplaneCLIVersion
}

func (c *HcConfig) GetOcmLongLivedTokenPath() string {
	return c.OcmLongLivedTokenPath
}
