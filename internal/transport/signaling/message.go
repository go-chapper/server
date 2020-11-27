// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities to handle WebRTC signaling
package signaling

type Message interface {
	Handle(*Hub) error
	Type() string
}
