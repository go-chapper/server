// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package installer provides ...
package installer

import (
	"strconv"
	"time"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/store"
)

func (i *Installer) defaultConfiguration() error {
	i.spinner.Suffix = " Creating default configuration"
	i.spinner.Start()

	cfg := config.NewDefault()
	cfg.General.Name = i.answers.InstanceName
	cfg.Router.Domain = i.answers.InstanceDomain

	if i.Type == "Advanced" {
		port, err := strconv.Atoi(i.answers.DatabasePort)
		if err != nil {
			return err
		}

		cfg.Store = config.StoreOptions{
			User:     i.answers.DatabaseUser,
			Password: i.answers.DatabasePassword,
			Database: i.answers.DatabaseName,
			Host:     i.answers.DatabaseHost,
			Port:     port,
		}
	}

	i.config = cfg
	time.Sleep(time.Second)
	return nil
}

func (i *Installer) databaseConnection() error {
	i.spinner.Suffix = " Connecting to database"
	s, err := store.New(i.answers.DatabaseType, i.config.Store)
	if err != nil {
		return err
	}

	i.store = s
	time.Sleep(time.Second)
	return nil
}

func (i *Installer) databaseMigration() error {
	i.spinner.Suffix = " Migrating the database"
	err := i.store.Migrate()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return nil
}

func (i *Installer) writeConfiguration() error {
	i.spinner.Suffix = " Saving configuration file"
	// TODO <2020/09/09>: Don't hardcode this path
	err := i.config.Write("chapper.toml")
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	i.spinner.Stop()
	return nil
}
