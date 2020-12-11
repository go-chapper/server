// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

type RoomService struct {
	store *store.Store
}

func NewRoomService(store *store.Store) RoomService {
	return RoomService{
		store: store,
	}
}

func (s RoomService) CreateRoom(room *models.Room) error {
	room.Hash = hash.FNV64(room.Name)
	return s.store.CreateRoom(room)
}

func (s RoomService) GetRoom(serverHash string) (models.Room, error) {
	return s.store.GetRoom(serverHash)
}

func (s RoomService) GetRooms() ([]models.Room, error) {
	return s.store.GetRooms()
}

func (s RoomService) UpdateRoom(roomHash string) error {
	return nil
}

func (s RoomService) DeleteRoom(roomHash string) error {
	return s.store.DeleteRoom(roomHash)
}
