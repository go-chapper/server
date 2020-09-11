// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cmd provides command line commands
package cmd

import (
	"fmt"
	"log"

	"git.web-warrior.de/go-chapper/server/internal/config"
	"git.web-warrior.de/go-chapper/server/internal/constants"
	"git.web-warrior.de/go-chapper/server/internal/logger"
	"git.web-warrior.de/go-chapper/server/internal/router"
	"git.web-warrior.de/go-chapper/server/internal/router/handlers"
	"git.web-warrior.de/go-chapper/server/internal/store"

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

		c := config.New()
		err := c.Read(configFilePath)
		if err != nil {
			fmt.Printf("ERROR [Config] Could not read config file: %s\n", err.Error())
		}

		f, err := logger.New(c.Log)
		if err != nil {
			fmt.Printf("ERROR [Logger] Failed to open log file: %s\n", err.Error())
		}
		defer f.Close()

		s, err := store.New("mysql", c.Store)
		if err != nil {
			log.Fatalf("ERROR [Store] Failed to connect to store: %v\n", err)
		}

		err = s.Migrate()
		if err != nil {
			log.Fatalf("ERROR [Store] Failed to migrate: %v\n", err)
		}

		r := router.New(c)
		h := handlers.New(s, c)
		r.AddRoutes(h)

		err = r.Run()
		if err != nil {
			log.Fatalf("ERROR [Router] Failed to start the router: %v\n", err)
		}
	},
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "Path to your config file")
	cobra.MarkFlagRequired(runCmd.PersistentFlags(), "config")

	rootCmd.AddCommand(runCmd)
}
