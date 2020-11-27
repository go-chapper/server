// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package call provides utilities to create, join and leave calls
package call

import (
	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/transport/bridge"
)

type Service struct {
	hub *bridge.Hub
}

func NewService(turn config.TurnOptions) Service {
	return Service{
		hub: bridge.NewHub(turn),
	}
}

func (s Service) NewCall(room string) error {
	return s.hub.NewCall(room)
}

func (s Service) ForwardSDP(roomHash, sdp string) error {
	return s.hub.ForwardSDP(roomHash, sdp)
}
