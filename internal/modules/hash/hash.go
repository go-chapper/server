// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hash provides multiple hash functions
package hash

type Hash interface {
	Hash(string) (string, error)
	Valid(string, string) (bool, error)
}
