// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

// GetSignalingChannel opens a websocket used for signaling
// func (h *Handler) GetSignalingChannel(c echo.Context) error {
// 	conn, err := h.signalingHub.CreateConnection(c.Response(), c.Request())
// 	if err != nil {
// 		log.Printf("ERROR [Router] Failed to upgrade connection: %v\n", err)
// 		return err
// 	}

// 	go h.signalingHub.Register(conn)
// 	go conn.ListenRead()
// 	go conn.ListenWrite()

// 	return nil
// }

// // GetSignalingToken returns an auth token to subscribe to the signaling websocket
// func (h *Handler) GetSignalingToken(c echo.Context) error {
// 	claims := getClaimes(c)

// 	// TODO <2020/13/09>: Add check if callee is blocked and/or the caller is blocked by callee

// 	t, err := h.signalingHub.Token(claims.Username)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, Map{
// 			"error": ErrInternal,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, Map{
// 		"token": t,
// 	})
// }
