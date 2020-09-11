// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package installer provides utilities for installing a Chapper instance
package installer

import "github.com/AlecAivazis/survey/v2"

func QuickSteps() []*survey.Question {
	return []*survey.Question{
		{
			Name: "instance-name",
			Prompt: &survey.Input{
				Message: "Instance | Name:",
				Default: "Chapper",
			},
		},
		{
			Name: "instance-domain",
			Prompt: &survey.Input{
				Message: "Instance | Domain (e.g. chapper.example.com):",
				Default: "",
			},
		},
		{
			Name: "database-type",
			Prompt: &survey.Select{
				Message: "Database | Type",
				Options: []string{"MySQL"},
				Default: "MySQL",
			},
		},
		{
			Name: "admin-panel-password",
			Prompt: &survey.Password{
				Message: "Please enter a password to protect the admin panel:",
			},
		},
	}
}
