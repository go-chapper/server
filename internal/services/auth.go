// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"net/http"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/avatar"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/modules/twofa"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
)

var (
	ErrCreateVerifyToken = NewError("create-verify-token", "Failed to create 2fa verify token", http.StatusInternalServerError)
	ErrUpdateVerifyToken = NewError("update-verify-token", "Failed to update 2fa verify token", http.StatusInternalServerError)
	ErrMissingDataUser   = NewError("missing-data-user", "Some data to login/register is missing", http.StatusBadRequest)
	ErrInvalidPassword   = NewError("invalid-password", "The user provided an invalid password", http.StatusUnauthorized)
	ErrBindUser          = NewError("bind-user", "Tailed to bind the user model", http.StatusInternalServerError)
	ErrCreateAvatar      = NewError("create-avatar", "Failed to create avatar", http.StatusInternalServerError)
	ErrHashPassword      = NewError("hash-password", "Failed to hash password", http.StatusInternalServerError)
	ErrSignToken         = NewError("sign-token", "Failed to sign jwt token", http.StatusInternalServerError)
	ErrCreateUser        = NewError("create-user", "Failed to create user", http.StatusInternalServerError)
	ErrGetUser           = NewError("get-user", "Failed to get user", http.StatusInternalServerError)
)

// AuthService wraps authentication dependencies
type AuthService struct {
	hash   hash.Hash
	store  *store.Store
	config *config.Config
	logger *log.Logger
}

// NewAuthService returns a new authentication service
func NewAuthService(store *store.Store, config *config.Config, logger *log.Logger) AuthService {
	return AuthService{
		hash:   hash.NewArgon2(),
		store:  store,
		config: config,
		logger: logger,
	}
}

// Register handles the registration process of a new user
func (s AuthService) Register(c echo.Context) error {
	var user models.PublicUser

	// Bind to signup user model
	err := c.Bind(&user)
	if err != nil {
		s.logger.Error(err)
		return ErrBindUser
	}

	// Check if some data is missing
	if user.IsEmpty() {
		s.logger.Info("Auth service: Some data to login/register is missing")
		return ErrMissingDataUser
	}

	// Hach the password to save into the database
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		s.logger.Error(err)
		return ErrHashPassword
	}
	user.Password = hashedPassword

	// Insert new user into the database
	err = s.store.CreateUser(user)
	if err != nil {
		s.logger.Error(err)
		return ErrCreateUser
	}

	// Generate the default profile avatar based on the username
	profileAvatar := avatar.New(240, user.Username)
	err = profileAvatar.Generate(s.config.Router.AvatarPath)
	if err != nil {
		s.logger.Error(err)
		return ErrCreateAvatar
	}

	return nil
}

// Login handles the login process of a user
func (s AuthService) Login(c echo.Context) (string, error) {
	var user models.PublicUser

	// Bind to user model
	err := c.Bind(&user)
	if err != nil {
		s.logger.Error(err)
		return "", ErrBindUser
	}

	// Check if some data is missing
	if user.IsLoginEmpty() {
		s.logger.Info("Auth service: Some data to login/register is missing")
		return "", ErrMissingDataUser
	}

	// Get the account from the database by username
	account, err := s.store.GetUser(user.Username)
	if err != nil {
		s.logger.Error(err)
		return "", ErrGetUser
	}

	// Check if the account uses 2FA
	if account.UsesTwoFA() {
		verifyToken, err := s.GenerateVerifyToken()
		if err != nil {
			s.logger.Error(err)
			return "", ErrCreateVerifyToken
		}

		err = s.store.UpdateTwoFAVerify(account.Username, verifyToken)
		if err != nil {
			s.logger.Error(err)
			return "", ErrUpdateVerifyToken
		}

		return "", nil
	}

	// Compare the provided with the saved password
	valid, err := s.ComparePassword(user.Password, account.Password)
	if !valid || err != nil {
		s.logger.Error(err)
		return "", ErrInvalidPassword
	}

	// Generate a new JWT token
	token := s.NewJWT(s.config.Router.JWTSecret, &jwt.Claims{
		Username: account.Username,
	})

	signedToken, err := token.Sign()
	if err != nil {
		s.logger.Error(err)
		return "", ErrSignToken
	}

	return signedToken, nil
}

// HashPassword returns the Argon2-hashed password or an error
func (s AuthService) HashPassword(password string) (string, error) {
	h, err := s.hash.Hash(password)
	if err != nil {
		return "", err
	}

	return h, nil
}

// ComparePassword returns if the input and the real password match
func (s AuthService) ComparePassword(input, hashed string) (bool, error) {
	return s.hash.Valid(input, hashed)
}

// NewJWT returns a new JWT
func (s AuthService) NewJWT(secret string, c *jwt.Claims) *jwt.JWT {
	return jwt.New(secret, c)
}

// ParseJWT parses a JWT
func (s AuthService) ParseJWT(input, key string, c jwt.Claims) (*jwt.JWT, error) {
	return jwt.Parse(input, key, c)
}

// GenerateTOTP generates a new TOTP
func (s AuthService) GenerateTOTP(issuer, account string) (twofa.TOTPKey, error) {
	return twofa.GenerateTOTP(issuer, account)
}

// ValidateTOTP validates a 6 digit TOTP code
func (s AuthService) ValidateTOTP(code, secret string) bool {
	return twofa.ValidateTOTP(code, secret)
}

// GenerateVerifyToken generates a random verify token
func (s AuthService) GenerateVerifyToken() (string, error) {
	return utils.RandomCryptoString(16)
}
