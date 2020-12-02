// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewCall initiates a new call
func (h *Handler) NewCall(c echo.Context) error {
	// TODO <2020/25/11>: Add room type check

	fmt.Println("Hey there")

	// err := h.callService.NewCall(c.Param("room-hash"))
	// if err != nil {
	// 	log.Printf("ERROR [Router] Failed to start new call: %v\n", err)
	// 	return c.JSON(http.StatusInternalServerError, Map{
	// 		"errror": ErrInternal,
	// 	})
	// }

	return c.JSON(http.StatusOK, Map{
		"status": "call-created",
	})
}

func (h *Handler) JoinCall(c echo.Context) error {
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
