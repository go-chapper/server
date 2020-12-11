// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cmd provides command line commands
package cmd

import (
	"fmt"
	"os"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/constants"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/router"
	"chapper.dev/server/internal/router/handlers"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/transport/turn"

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
			fmt.Printf("[E] config: could not read config file: %v\n", err)
			os.Exit(1)
		}

		logger, err := log.New(c.Log)
		if err != nil {
			fmt.Printf("[E] logger: failed to setup logger: %v\n", err)
			os.Exit(1)
		}

		s, err := store.New("mysql", c.Store)
		if err != nil {
			e := fmt.Errorf("store: failed to connect to database: %v", err)
			logger.Fatal(e)
		}

		err = s.Migrate()
		if err != nil {
			e := fmt.Errorf("store: failed to migrate table(s): %v", err)
			logger.Fatal(e)
		}

		r := router.New(c)
		h := handlers.New(s, c)
		r.AddRoutes(h)

		t, err := turn.New(c.Turn.PublicIP, c.Router.Domain, "udp4", c.Turn.Port)
		if err != nil {
			e := fmt.Errorf("turn: failed to init server: %v", err)
			logger.Fatal(e)
		}

		err = t.Run()
		if err != nil {
			e := fmt.Errorf("turn: failed to start server: %v", err)
			logger.Fatal(e)
		}

		err = r.Run()
		if err != nil {
			e := fmt.Errorf("router: failed to start router: %v", err)
			logger.Fatal(e)
		}
	},
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "Path to your config file")
	cobra.MarkFlagRequired(runCmd.PersistentFlags(), "config")

	rootCmd.AddCommand(runCmd)
}
