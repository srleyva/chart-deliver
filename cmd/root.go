// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/srleyva/chart-deliver/pkg/helpers"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var meta helpers.Template

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chart",
	Short: "Helm Chart Dynamic Generation",
	Long: `This is designed to be used as a stage of a pipeline.
The idea being that a project can provide a values.yaml
file and a chart will be dynamically generated and deployed`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Lookup("name").Value.String() == "" || cmd.Flags().Lookup("release").Value.String() == "" {
			return fmt.Errorf("supply name and release name")
		}
		fflags := cmd.Flags()
		runner := helpers.NewHelmHandler()
		meta = helpers.Template{
			Runner:      runner,
			ChartName:   fflags.Lookup("name").Value.String(),
			ReleaseName: fflags.Lookup("release").Value.String(),
			Version:     fflags.Lookup("version").Value.String(),
			Values:      fflags.Lookup("values").Value.String(),
			Image:       fflags.Lookup("image").Value.String(),
			Tag:         fflags.Lookup("tag").Value.String(),
			Namespace:   fflags.Lookup("namespace").Value.String(),
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Info("Chart Deliver Version v0.0.1")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("err: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFormatter(&log.JSONFormatter{})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chart-deliver.yaml)")
	rootCmd.PersistentFlags().StringP("release", "r", "", "The Release name of the deploy")
	rootCmd.PersistentFlags().StringP("name", "n", "", "The name of the generated helm chart")
	rootCmd.PersistentFlags().StringP("version", "v", "v0.0.1", "The version of the generated helm chart")
	rootCmd.PersistentFlags().StringP("values", "f", "", "The path to the file containing your values")
	rootCmd.PersistentFlags().StringP("image", "c", "", "The path to the file containing your values")
	rootCmd.PersistentFlags().StringP("tag", "t", "", "The path to the file containing your values")

	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("release", rootCmd.PersistentFlags().Lookup("release"))
	viper.BindPFlag("values", rootCmd.PersistentFlags().Lookup("values"))
	viper.BindPFlag("image", rootCmd.PersistentFlags().Lookup("image"))
	viper.BindPFlag("tag", rootCmd.PersistentFlags().Lookup("tag"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".chart-deliver" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".chart-deliver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
