// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package call provides utilities to create, join and leave calls
package call

import (
	"net/http"

	"chapper.dev/server/internal/transport/bridge"
)

type Service struct {
	bridge *bridge.Bridge
}

func NewService() Service {
	return Service{
		bridge: bridge.NewBridge(),
	}
}

func (s Service) NewCall(username, roomHash string, w http.ResponseWriter, r *http.Request) error {
	return s.bridge.Connect(username, roomHash, w, r)
}
