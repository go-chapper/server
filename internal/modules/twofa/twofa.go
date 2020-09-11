// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package twofa provides utilities to generate and validate 2FA codes
package twofa

import (
	"image"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPKey struct {
	key *otp.Key
}

// GenerateTOTP generates a new TOTP
func GenerateTOTP(issuer, account string) (TOTPKey, error) {
	options := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
		// NOTE(Techassi): Use SHA1 because Google Authenticator does BS when using SHA512
		// Maybe make this algo configurable by the user
		Algorithm: otp.AlgorithmSHA1,
	}

	key, err := totp.Generate(options)
	return TOTPKey{
		key: key,
	}, err
}

// ValidateTOTP validates a 6 digit TOTP code
func ValidateTOTP(code, secret string) bool {
	ok, _ := totp.ValidateCustom(
		code,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
	return ok
}

// Secret returns the TOTP secret
func (k TOTPKey) Secret() string {
	return k.key.Secret()
}

// Image returns the TOTP as an image to add to 2FA Managers
func (k TOTPKey) Image(width, height int) (image.Image, error) {
	return k.key.Image(width, height)
}

// URL returns the TOTP URL
func (k TOTPKey) URL() string {
	return k.key.URL()
}
