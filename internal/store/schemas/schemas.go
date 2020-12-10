// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package schemas provides database schemas
package schemas

func All() []string {
	return []string{Users, Servers, Rooms, Invites}
}
