// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"git.web-warrior.de/go-chapper/server/internal/models"
)

// CreateRoom creates a new 'room' on server identified by 'serverhash'
func (s *Store) CreateRoom(serverhash string, room *models.Room) error {
	return s.Ctx().Model(&models.Server{Hash: serverhash}).Association("Rooms").Append(room)
}

// GetRoom returns one room of a server identified by 'serverhash'
func (s *Store) GetRoom(serverhash, hash string) (models.Room, error) {
	var room models.Room
	return room, s.Ctx().Model(&models.Server{Hash: serverhash}).Association("Rooms").Find(&room, "hash = ?", hash)
}

// GetRooms returns all rooms of a server identified by 'serverhash'
func (s *Store) GetRooms(serverhash string) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, s.Ctx().Model(&models.Server{Hash: serverhash}).Association("Rooms").Find(&rooms)
}
