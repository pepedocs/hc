package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkgInt "hc/internal"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "hc",
	Short: "Hybric cloud containerized environment",
	Long: `A CLI for locally provisioning a container that provides a management 
	environment for managing an OpenShift-based hybrid cloud.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&pkgInt.Debug, "debug", "d", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hc.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		var src string
		if os.Getenv("IS_IN_CONTAINER") == "true" {
			src = "/"
		} else {
			home, err := os.UserHomeDir()
			src = home
			cobra.CheckErr(err)
		}
		viper.AddConfigPath(src)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	if err = pkgInt.ValidateConfig(); err != nil {
		log.Fatalf("Invalid config: %s", err)
	}

	// Todo: Handle os signals

}
