// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"log"

	"github.com/labstack/echo/v4"
)

// NewCall initiates a new call
func (h *Handler) NewCall(c echo.Context) error {
	claims := getClaimes(c)
	roomHash := c.Param("room-hash")

	err := h.callService.NewCall(claims.Username, roomHash, c.Response().Writer, c.Request())
	if err != nil {
		log.Printf("ERROR [Router] Unable to create new call: %v\n", err)
		return err
	}

	return nil
}

func (h *Handler) JoinCall(c echo.Context) error {
	// claims := getClaimes(c)
	roomHash := c.Param("room-hash")

	err := h.callService.NewCall("Test", roomHash, c.Response().Writer, c.Request())
	if err != nil {
		log.Printf("ERROR [Router] Unable to create or join call: %v\n", err)
		return err
	}

	return nil
}

func (h *Handler) ForwardSDP(c echo.Context) error {
	// sdp, err := ioutil.ReadAll(c.Request().Body)
	// if err != nil {
	// 	log.Printf("ERROR [Router] Failed read sdp: %v\n", err)
	// 	return c.JSON(http.StatusInternalServerError, Map{
	// 		"errror": ErrInternal,
	// 	})
	// }

	// err = h.callService.ForwardSDP(c.Param("room-hash"), string(sdp))
	// if err != nil {
	// 	log.Printf("ERROR [Router] Failed forward sdp: %v\n", err)
	// 	return c.JSON(http.StatusInternalServerError, Map{
	// 		"errror": ErrInternal,
	// 	})
	// }

	return nil
}
