// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package direct provides utilities to create, get, update and delete direct chats/calls
package direct

import "chapper.dev/server/internal/store"

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return Service{
		store: store,
	}
}
