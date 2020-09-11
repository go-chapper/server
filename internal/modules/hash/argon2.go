// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hash provides multiple hash functions
package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"chapper.dev/server/internal/constants"
	"chapper.dev/server/internal/utils"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("The encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("Incompatible version of argon2")
)

type Argon2 struct {
	name   string
	config Argon2Config
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2Config is the default config when using NewArgon2
var DefaultArgon2Config = Argon2Config{
	Memory:      64 * constants.Argon2Kibibyte,
	Iterations:  5, // NOTE(Techassi): The RFC suggests 1 iteration, should we use 5?
	Parallelism: uint8(runtime.NumCPU() / 4),
	SaltLength:  16,
	KeyLength:   32,
}

// NewArgon2 returns a new Argon2 hasher with the default config
func NewArgon2() Argon2 {
	return NewArgon2WithConfig(DefaultArgon2Config)
}

// NewArgon2WithConfig returns a new Argon2 hasher with a custom config
func NewArgon2WithConfig(c Argon2Config) Argon2 {
	return Argon2{
		name:   "argon2",
		config: c,
	}
}

// Name returns the name of the hasher
func (a Argon2) Name() string {
	return a.name
}

// Hash hashes the payload
func (a Argon2) Hash(payload string) (string, error) {
	// Generate cryptographically secure random salt
	salt, err := utils.RandomByteSlice(a.config.SaltLength)
	if err != nil {
		return "", err
	}

	// Generate the salted hash of the password
	hash := argon2.IDKey([]byte(payload),
		salt,
		a.config.Iterations,
		a.config.Memory,
		a.config.Parallelism,
		a.config.KeyLength,
	)

	// Generate the encoded representation of the password
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	passB64 := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.config.Memory,
		a.config.Iterations,
		a.config.Parallelism,
		saltB64,
		passB64,
	)
	return encoded, nil
}

// Valid compares the input to a hashed payload
func (a Argon2) Valid(input, hashed string) (bool, error) {
	p, salt, hash, err := decodeHash(hashed)
	if err != nil {
		return false, err
	}

	compareHash := argon2.IDKey([]byte(input), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	if subtle.ConstantTimeCompare(hash, compareHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(encodedHash string) (p *Argon2Config, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Argon2Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
