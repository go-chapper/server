// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"chapper.dev/server/internal/models"
)

// CreateServer inserts a new server entry into the database
func (s *Store) CreateServer(server *models.Server) error {
	_, err := s.conn.Exec(`
		INSERT INTO servers
		(hash, name, description, image)
		VALUES (?, ?, ?, ?)`,
		server.Hash,
		server.Name,
		server.Description,
		server.Image,
	)
	return err
}

// GetServer selects ONE server entry with provided 'serverHash' from the database
func (s *Store) GetServer(serverHash string) (models.Server, error) {
	var server models.Server
	err := s.conn.Get(&server,
		`SELECT hash, name, description, image
		FROM servers
		WHERE hash = ?`,
		serverHash,
	)
	return server, err
}

// GetServers selects multiple server entries from the database
func (s *Store) GetServers() ([]models.Server, error) {
	var servers []models.Server
	err := s.conn.Select(&servers, `SELECT hash, name, description, imageFROM servers`)
	return servers, err
}

// UpdateServer updates ONE server entry with provided 'serverHash' in the database
func (s *Store) UpdateServer(serverHash string, new *models.Server) error {
	_, err := s.conn.Exec(`
		UPDATE servers
		SET name = ?, description = ?, image = ?
		WHERE hash = ?`,
		new.Name,
		new.Description,
		new.Image,
		serverHash,
	)
	return err
}

// DeleteServer deletes ONE server entry with provided 'serverHash' from the database
func (s *Store) DeleteServer(serverHash string) error {
	_, err := s.conn.Exec(`
		DELETE FROM servers
		WHERE hash = ?`,
		serverHash,
	)
	return err
}
