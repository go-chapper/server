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

func (s Service) CreateRoom(room *models.Room) error {
	room.Hash = hash.FNV64(room.Name)
	return s.store.CreateRoom(room)
}

func (s Service) GetRoom(serverHash string) (models.Room, error) {
	return s.store.GetRoom(serverHash)
}

func (s Service) GetRooms() ([]models.Room, error) {
	return s.store.GetRooms()
}

func (s Service) UpdateRoom(roomHash string) error {
	return nil
}

func (s Service) DeleteRoom(roomHash string) error {
	return s.store.DeleteRoom(roomHash)
}
