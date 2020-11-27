// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hash

import (
	"encoding/hex"
	"hash/fnv"
)

// FNV64 returns the FNV64 hash of 'payload'
func FNV64(payload string) string {
	h := fnv.New64()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// FNV128 returns the FNV128 hash of 'payload'
func FNV128(payload string) string {
	h := fnv.New128()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
