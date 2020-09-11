// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package installer provides utilities for installing a Chapper instance
package installer

import (
	"fmt"
	"time"

	"git.web-warrior.de/go-chapper/server/internal/config"
	"git.web-warrior.de/go-chapper/server/internal/constants"
	"git.web-warrior.de/go-chapper/server/internal/store"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
)

type Installer struct {
	Type     string
	Step     int
	quick    []*survey.Question
	advanced []*survey.Question
	answers  *Answers
	spinner  *spinner.Spinner
	config   *config.Config
	store    *store.Store
}

type Answers struct {
	InstanceName       string `survey:"instance-name"`
	InstanceDomain     string `survey:"instance-domain"`
	DatabaseType       string `survey:"database-type"`
	DatabaseUser       string `survey:"database-user"`
	DatabasePassword   string `survey:"database-password"`
	DatabaseName       string `survey:"database-name"`
	DatabaseHost       string `survey:"database-host"`
	DatabasePort       string `survey:"database-port"`
	FrontendAuto       bool   `survey:"frontend-auto"`
	FrontendPath       string `survey:"frontend-path"`
	ConfigPath         string `survey:"config-path"`
	AdminPanelPassword string `survey:"admin-panel-password"`
}

func New() *Installer {
	return &Installer{
		Step:     0,
		quick:    QuickSteps(),
		advanced: AdvancedSteps(),
		answers:  new(Answers),
	}
}

func (i *Installer) Start() error {
	var installMethod = ""
	err := survey.AskOne(&survey.Select{
		Message: "Installation method:",
		Options: []string{"Quick", "Advanced"},
	}, &installMethod)
	if err != nil {
		return err
	}

	i.Type = installMethod
	if installMethod == "Quick" {
		err := survey.Ask(i.quick, i.answers)
		if err != nil {
			return err
		}
	} else {
		err := survey.Ask(i.advanced, i.answers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Installer) GetAnswers() *Answers {
	return i.answers
}

func (i *Installer) InstallServer() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.FinalMSG = fmt.Sprintf("%s>%s Complete. Happy chatting!\n", constants.ColorGreen, constants.ColorReset)
	s.Color("green")
	i.spinner = s

	err := i.defaultConfiguration()
	if err != nil {
		return err
	}

	err = i.databaseConnection()
	if err != nil {
		return err
	}

	err = i.databaseMigration()
	if err != nil {
		return err
	}

	return i.writeConfiguration()
}
