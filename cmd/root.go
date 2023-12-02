/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkgInt "github.com/hc/internal"
)

var cfgFile string
var config *pkgInt.HcConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hc",
	Short: "Hybrid cloud container",
	Long: `A CLI locally provisioning a container which provides a management 
	environment for managing an OpenShift-based hybrid cloud.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
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
		viper.SetConfigName(".ocm-workspace")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		glog.Fatal("Failed to read config file: ", err)
	}

	config = pkgInt.NewHcConfig()
}
