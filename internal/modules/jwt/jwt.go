// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package jwt provides utilities to genearte and validate JWT credentials
package jwt

import (
	"errors"

	"chapper.dev/server/internal/models"
	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrUsernameEmpty indicates that username is empty in claims
	ErrUsernameEmpty = errors.New("Username cannot be empty")
)

// JWT wraps a JWT token, it's key and claims
type JWT struct {
	token  *jwt.Token
	key    string
	claims Claims
}

// Claims is a custom claims struct
type Claims struct {
	Username   string            `json:"username"`
	Privileges models.Privileges `json:"privileges"`
	StandardClaims
}

// StandardClaims is a wrapper for jwt.StandardClaims
type StandardClaims jwt.StandardClaims

// New returns a new JWT instance
func New(key string, c *Claims) *JWT {
	return &JWT{
		token: jwt.NewWithClaims(jwt.SigningMethodHS256, c),
		key:   key,
	}
}

// Parse parses the input token and returns a new JWT instance or an error
func Parse(input, key string, c Claims) (*JWT, error) {
	token, err := jwt.ParseWithClaims(input, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	return &JWT{
		token:  token,
		key:    key,
		claims: c,
	}, nil
}

// Claims returns the cliams of the token
func (j *JWT) Claims() Claims {
	return j.claims
}

// Sign signs the token and returns the token as a string or an error
func (j *JWT) Sign() (string, error) {
	return j.token.SignedString([]byte(j.key))
}

// Valid returns wether the claims are valid
func (c Claims) Valid() error {
	if c.Username == "" {
		return ErrUsernameEmpty
	}
	return nil
}
