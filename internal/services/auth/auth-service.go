// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package auth provides utilities used to authenticte users
package auth

import (
	"log"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/avatar"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/modules/twofa"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	ErrBindUser          = errors.New("bind-user")
	ErrMissingDataUser   = errors.New("missing-data-user")
	ErrHashPassword      = errors.New("hash-password")
	ErrCreateUser        = errors.New("create-user")
	ErrGetUser           = errors.New("get-user")
	ErrCreateAvatar      = errors.New("create-avatar")
	ErrCreateVerifyToken = errors.New("create-verify-token")
	ErrUpdateVerifyToken = errors.New("update-verify-token")
	ErrInvalidPassword   = errors.New("invalid-password")
)

// Service wraps authentication dependencies
type Service struct {
	hash   hash.Hash
	store  *store.Store
	config *config.Config
	// logger *logger.Logger
}

// NewService returns a new authentication service
func NewService(store *store.Store, config *config.Config) Service {
	return Service{
		hash:   hash.NewArgon2(),
		store:  store,
		config: config,
	}
}

// Register handles the registration process of a new user
func (s Service) Register(c echo.Context) error {
	var user models.SignupUser

	// Bind to signup user model
	err := c.Bind(&user)
	if err != nil {
		log.Printf("[E] auth service: %v\n", err)
		return ErrBindUser
	}

	// Check if some data is missing
	if user.IsEmpty() {
		return ErrMissingDataUser
	}

	// Hach the password to save into the database
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return ErrHashPassword
	}
	user.Password = hashedPassword

	// Insert new user into the database
	err = s.store.CreateUser(user)
	if err != nil {
		return ErrCreateUser
	}

	// Generate the default profile avatar based on the username
	profileAvatar := avatar.New(240, user.Username)
	err = profileAvatar.Generate(s.config.Router.AvatarPath)
	if err != nil {
		return ErrCreateAvatar
	}

	return nil
}

// Login handles the login process of a user
func (s Service) Login(c echo.Context) (string, error) {
	var user models.User

	// Bind to user model
	err := c.Bind(&user)
	if err != nil {
		// TODO <2020/11/12>: Log
		return "", ErrBindUser
	}

	// Check if some data is missing
	if user.IsLoginEmpty() {
		return "", ErrMissingDataUser
	}

	// Get the account from the database by username
	account, err := s.store.GetUser(user.Username)
	if err != nil {
		return "", ErrGetUser
	}

	// Check if the account uses 2FA
	if account.UsesTwoFA() {
		verifyToken, err := s.GenerateVerifyToken()
		if err != nil {
			return "", ErrCreateVerifyToken
		}

		err = s.store.UpdateTwoFAVerify(account.Username, verifyToken)
		if err != nil {
			return "", ErrUpdateVerifyToken
		}

		return "", nil
	}

	// Compare the provided with the saved password
	valid, err := s.ComparePassword(user.Password, account.Password)
	if !valid || err != nil {
		return "", ErrInvalidPassword
	}

	// Generate a new JWT token
	token := s.NewJWT(s.config.Router.JWTSecret, &jwt.Claims{
		Username: account.Username,
	})

	signedToken, err := token.Sign()
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// HashPassword returns the Argon2-hashed password or an error
func (s Service) HashPassword(password string) (string, error) {
	h, err := s.hash.Hash(password)
	if err != nil {
		return "", err
	}

	return h, nil
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

// GenerateVerifyToken generates a random verify token
func (s Service) GenerateVerifyToken() (string, error) {
	return utils.RandomCryptoString(16)
}
