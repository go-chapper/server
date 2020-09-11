// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cmd provides command line commands
package cmd

import (
	"fmt"
	"os"

	"chapper.dev/server/internal/constants"
	"chapper.dev/server/internal/installer"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Launches an interactive CLI to help you setting up your Chapper instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(constants.Banner, Version, constants.Website)
		fmt.Printf("\nWelcome to the interactive install CLI\n")
		fmt.Printf("--------------------------------------\n")

		// Start the installer
		ins := installer.New()
		err := ins.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Install / setup the server
		err = ins.InstallServer()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printNextSteps(ins.GetAnswers())
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func printNextSteps(answers *installer.Answers) {
	fmt.Printf("\nNext steps\n")
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Go ahead and run\n\n")
	fmt.Printf("    %schapper run%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %schapper run --config path/to/config%s\n\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("Visit https://%s/register to register your new account!\n", answers.InstanceDomain)
}
