// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package joins provides SQL join specific models
package joins

type UserServers struct {
	Hash        string `json:"hash"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"string"`
	Image       string `json:"image"`
}
