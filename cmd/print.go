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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srleyva/chart-deliver/pkg/helpers"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		runner := helpers.NewHelmHandler()
		meta := helpers.Template{
			Runner:      runner,
			ChartName:   fflags.Lookup("name").Value.String(),
			ReleaseName: fflags.Lookup("release").Value.String(),
			Version:     fflags.Lookup("version").Value.String(),
			Values:      fflags.Lookup("values").Value.String(),
		}

		log.Infof("Generating helm chart: %s", meta.ChartName)
		log.Infof("ReleaseName: %s", meta.ReleaseName)
		if err := meta.GenerateHelmChart(); err != nil {
			log.Fatalf("err generating file: %s", err)
		}
		log.Info("Chart Successfully Generated")
		log.Info("Printing Kubernetes spec")
		fmt.Print(meta.PrintHelmTemplate())
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
