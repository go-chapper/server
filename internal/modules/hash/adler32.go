// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hash provides multiple hash functions
package hash

import (
	"encoding/hex"
	"hash/adler32"
)

// Adler32 returns the Adler32 hash of 'payload'
func Adler32(payload string) string {
	h := adler32.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
