// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"chapper.dev/server/internal/app"
	"chapper.dev/server/internal/constants"

	"github.com/spf13/cobra"
)

var configFilePath string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run runs the Chapper server instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(constants.Banner, Version, constants.Website)

		a, err := app.New(configFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = a.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = a.Stop(ctx)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "Path to your config file")
	cobra.MarkFlagRequired(runCmd.PersistentFlags(), "config")

	rootCmd.AddCommand(runCmd)
}
