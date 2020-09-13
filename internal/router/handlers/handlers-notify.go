// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// GetNotifyChannel returns a new websocket
func (h *Handler) GetNotifyChannel(c echo.Context) error {
	conn, err := h.hub.CreateConnection(c.Response(), c.Request())
	if err != nil {
		log.Printf("ERROR [Router] Failed to upgrade connection: %v\n", err)
	}

	go h.hub.Register(conn)
	go conn.ListenRead()
	go conn.ListenWrite()

	return nil
}

func (h *Handler) Notify(c echo.Context) error {
	return nil
}
