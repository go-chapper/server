// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/labstack/echo/v4"
)

// GetMessagingChannel opens a websocket used for messaging
func (h *Handler) GetMessagingChannel(c echo.Context) error {
	// conn, err := h.messagingHub.CreateConnection(c.Response(), c.Request())
	// if err != nil {
	// 	log.Printf("ERROR [Router] Failed to upgrade connection: %v\n", err)
	// 	return err
	// }

	// go h.messagingHub.Register(conn)
	// go conn.ListenRead()
	// go conn.ListenWrite()

	return nil
}

// GetMessagingToken returns an auth token to subscribe to the messaging websocket
func (h *Handler) GetMessagingToken(c echo.Context) error {
	// claims := getClaimes(c)

	// // TODO <2020/13/09>: Add check if callee is blocked and/or the caller is blocked by callee

	// t, err := h.messagingHub.Token(claims.Username)
	// if err != nil {
	// 	log.Printf("ERROR [Router] Failed to fetch auth token for messaging connection: %v\n", err)
	// 	return c.JSON(http.StatusInternalServerError, Map{
	// 		"error": ErrInternal,
	// 	})
	// }

	// return c.JSON(http.StatusOK, Map{
	// 	"token": t,
	// })
	return nil
}
