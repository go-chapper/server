// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

func (s *Store) GetUser(username string) (models.User, error) {
	var user models.User
	// TODO <2020/10/12>: Join permissions
	err := s.conn.Get(&user,
		`SELECT username, password, email
		FROM users
		WHERE username = ?`,
		username,
	)
	return user, err
}

func (s *Store) GetUserPublicKey(username string) (string, error) {
	var publicKey string
	err := s.conn.Get(&publicKey,
		`SELECT publickey
		FROM users
		WHERE username = ?`,
		username,
	)
	return publicKey, err
}

func (s *Store) GetUserServers(username string) error {
	// TODO <2020/10/12>: re-implement
	return nil
}

func (s *Store) CreateUser(user models.PublicUser) error {
	_, err := s.conn.Exec(`
		INSERT INTO users
		(username, password, email, publickey)
		VALUES (?, ?, ?, ?)`,
		user.Username,
		user.Password,
		user.Email,
		user.PublicKey,
	)
	return err
}

func (s *Store) UpdateUser(username string, user *models.User) error {
	_, err := s.conn.Exec(`
		UPDATE users
		SET password = ?, email = ?
		WHERE username = ?`,
		user.Password,
		user.Email,
		username,
	)
	return err
}

func (s *Store) UpdateTwoFAVerify(username, verify string) error {
	_, err := s.conn.Exec(`
		UPDATE users
		SET twofa_verify = ?
		WHERE username = ?`,
		verify,
		username,
	)
	return err
}
