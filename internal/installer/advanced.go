// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package installer

import "github.com/AlecAivazis/survey/v2"

func AdvancedSteps() []*survey.Question {
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
			Name: "database-user",
			Prompt: &survey.Input{
				Message: "Database | User:",
				Default: "",
			},
		},
		{
			Name: "database-password",
			Prompt: &survey.Password{
				Message: "Database | Password:",
			},
		},
		{
			Name: "database-name",
			Prompt: &survey.Input{
				Message: "Database | Name:",
				Default: "",
			},
		},
		{
			Name: "database-host",
			Prompt: &survey.Input{
				Message: "Database | Host:",
				Default: "",
			},
		},
		{
			Name: "database-port",
			Prompt: &survey.Input{
				Message: "Database | Port:",
				Default: "3306",
			},
		},
		{
			Name: "frontend-auto",
			Prompt: &survey.Confirm{
				Message: "Frontend | Do you want want to automatically set up the Chapper frontend",
				Default: false,
			},
		},
		{
			Name: "frontend-path",
			Prompt: &survey.Input{
				Message: "Frontend | Root path (e.g. /var/www/chapper):",
			},
		},
		{
			Name: "config-path",
			Prompt: &survey.Input{
				Message: "Config | Path (e.g. /etc/chapper/config.toml):",
			},
		},
	}
}
