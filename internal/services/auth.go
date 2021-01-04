// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/avatar"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/modules/twofa"
	"chapper.dev/server/internal/services/errors"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
)

var (
	authCtx = log.NewContext("auth-srv")
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
		s.logger.Errorc(authCtx, err)
		return errors.ErrBindUser
	}

	// Check if some data is missing
	if user.IsInvalid() {
		s.logger.Infoc(authCtx, "some data to login/register is missing")
		return errors.ErrMissingUserData
	}

	// Hash the password to save into the database
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return errors.ErrHashPassword
	}
	user.Password = hashedPassword

	// Insert new user into the database
	err = s.store.CreateUser(user)
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return errors.ErrCreateUser
	}

	// Generate the default profile avatar based on the username
	profileAvatar := avatar.New(240, user.Username)
	err = profileAvatar.Generate(s.config.Router.AvatarPath)
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return errors.ErrCreateAvatar
	}

	return nil
}

// Login handles the login process of a user
func (s AuthService) Login(c echo.Context) (string, error) {
	var user models.PublicUser

	// Bind to user model
	err := c.Bind(&user)
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return "", errors.ErrBindUser
	}

	// Check if some data is missing
	if user.IsLoginEmpty() {
		s.logger.Infoc(authCtx, "some data to login/register is missing")
		return "", errors.ErrMissingUserData
	}

	// Get the account from the database by username
	account, err := s.store.GetUser(user.Username)
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return "", errors.ErrGetUser
	}

	// Check if the account uses 2FA
	if account.UsesTwoFA() {
		verifyToken, err := s.GenerateVerifyToken()
		if err != nil {
			s.logger.Errorc(authCtx, err)
			return "", errors.ErrCreateVerifyToken
		}

		err = s.store.UpdateTwoFAVerify(account.Username, verifyToken)
		if err != nil {
			s.logger.Errorc(authCtx, err)
			return "", errors.ErrUpdateVerifyToken
		}

		return "", nil
	}

	// Compare the provided with the saved password
	valid, err := s.ComparePassword(user.Password, account.Password)
	if !valid || err != nil {
		s.logger.Errorc(authCtx, err)
		return "", errors.ErrInvalidPassword
	}

	// Generate a new JWT token
	token := s.NewJWT(s.config.Router.JWTSecret, &jwt.Claims{
		Username: account.Username,
	})

	signedToken, err := token.Sign()
	if err != nil {
		s.logger.Errorc(authCtx, err)
		return "", errors.ErrSignToken
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
