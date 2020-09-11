// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package room provides utilities to create, get, update and delete rooms on virtual
// servers
package room

import (
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return Service{
		store: store,
	}
}

func (s Service) CreateRoom(serverhash string, room *models.Room) error {
	room.Hash = hash.Adler32(room.Name)
	return s.store.CreateRoom(serverhash, room)
}

func (s Service) GetRoom(serverhash, hash string) (models.Room, error) {
	return s.store.GetRoom(serverhash, hash)
}

func (s Service) GetRooms(serverhash string) ([]models.Room, error) {
	return s.store.GetRooms(serverhash)
}
