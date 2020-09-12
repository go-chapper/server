// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"log"
	"net/http"
	"strings"

	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/jwt"

	"github.com/labstack/echo/v4"
)

// AuthRegister registers a new user
func (h *Handler) AuthRegister(c echo.Context) error {
	newUser := new(models.SignupUser)
	err := c.Bind(newUser)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrBind,
		})
	}

	if newUser.IsEmpty() {
		log.Println("WARNING [Router] Missing/empty data for registration")
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrEmptyData,
		})
	}

	err = h.authService.HashPassword(newUser)
	if err != nil {
		log.Printf("ERROR [Router] Failed to hash password: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrHashPassword,
		})
	}

	err = h.userService.CreateUser(newUser)
	if err != nil {
		log.Printf("ERROR [Router] Failed to create new user: %v\n", err)
		// TODO <2020/06/09>: Can we do this more cleanly?
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, Map{
				"error": ErrUsernameTaken,
			})
		}

		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrCreateUser,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": StatusRegistered,
	})
}

// AuthRefresh refreshes the JWT token
func (h *Handler) AuthRefresh(c echo.Context) error {
	refresh := new(models.RefreshToken)
	err := c.Bind(refresh)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrBind,
		})
	}

	if refresh.IsEmpty() {
		log.Println("WARNING [Router] Missing/empty data to refresh token")
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrEmptyData,
		})
	}

	return nil
}

// AuthLogin logs a user in
func (h *Handler) AuthLogin(c echo.Context) error {
	user := new(models.User)
	err := c.Bind(user)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": "Invalid request",
			"code":  ErrBind,
		})
	}

	if user.IsLoginEmpty() {
		log.Println("WARNING [Router] Missing/empty data for login")
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrEmptyData,
		})
	}

	acc, err := h.userService.GetUser(user.Username)
	if err != nil {
		log.Printf("WARNING [Router] This username %s does not exist\n", user.Username)
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrGetUser,
		})
	}

	if acc.UsesTwoFA() {
		temp, err := h.authService.GenerateTempToken()
		if err != nil {
			log.Printf("ERROR [Router] Failed to generate temp token for 2FA: %v\n", err)
			return c.JSON(http.StatusInternalServerError, Map{
				"error": ErrTwoFA,
			})
		}

		err = h.userService.SaveTempToken(acc.Username, temp)
		if err != nil {
			log.Printf("ERROR [Router] Failed to save temp token for 2FA: %v\n", err)
			return c.JSON(http.StatusInternalServerError, Map{
				"error": ErrTwoFA,
			})
		}

		return c.JSON(http.StatusOK, Map{
			"action": "code",
			"token":  temp,
		})
	}

	valid, err := h.authService.ComparePassword(user.Password, acc.Password)
	if !valid || err != nil {
		log.Printf("WARNING [Router] Unauthorized: %v\n", err)
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	token := h.authService.NewJWT(h.config.Router.JWTSigningKey, &jwt.Claims{
		Username:   acc.Username,
		Privileges: acc.Role[0].Privileges,
	})

	signed, err := token.Sign()
	if err != nil {
		log.Printf("ERROR [Router] Failed to sign JWT: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrJWT,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"action": "login",
		"token":  signed,
	})
}
