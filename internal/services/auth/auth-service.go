// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package auth provides utilities used to authenticte users
package auth

import (
	"git.web-warrior.de/go-chapper/server/internal/models"
	"git.web-warrior.de/go-chapper/server/internal/modules/hash"
	"git.web-warrior.de/go-chapper/server/internal/modules/jwt"
	"git.web-warrior.de/go-chapper/server/internal/modules/twofa"
	"git.web-warrior.de/go-chapper/server/internal/utils"
)

// Service wraps authentication dependencies
type Service struct {
	hash hash.Hash
}

// NewService returns a new authentication service
func NewService() Service {
	return Service{
		hash: hash.NewArgon2(),
	}
}

// HashPassword returns the Argon2-hashed password or an error
func (s Service) HashPassword(user *models.User) error {
	h, err := s.hash.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = h
	return nil
}

// ComparePassword returns if the input and the real password match
func (s Service) ComparePassword(input, hashed string) (bool, error) {
	return s.hash.Valid(input, hashed)
}

// NewJWT returns a new JWT
func (s Service) NewJWT(secret string, c *jwt.Claims) *jwt.JWT {
	return jwt.New(secret, c)
}

// ParseJWT parses a JWT
func (s Service) ParseJWT(input, key string, c jwt.Claims) (*jwt.JWT, error) {
	return jwt.Parse(input, key, c)
}

// GenerateTOTP generates a new TOTP
func (s Service) GenerateTOTP(issuer, account string) (twofa.TOTPKey, error) {
	return twofa.GenerateTOTP(issuer, account)
}

// ValidateTOTP validates a 6 digit TOTP code
func (s Service) ValidateTOTP(code, secret string) bool {
	return twofa.ValidateTOTP(code, secret)
}

// GenerateTempToken generates a random temporary token
func (s Service) GenerateTempToken() (string, error) {
	return utils.RandomCryptoString(16)
}
