// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"net/http"

	"chapper.dev/server/internal/transport/bridge"
)

type CallService struct {
	bridge *bridge.Bridge
}

func NewCallService() CallService {
	return CallService{
		bridge: bridge.NewBridge(),
	}
}

func (s CallService) NewCall(username, roomHash string, w http.ResponseWriter, r *http.Request) error {
	return s.bridge.Connect(username, roomHash, w, r)
}
