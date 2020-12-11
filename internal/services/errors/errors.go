// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package errors provides service level errors
package errors

import (
	"fmt"
	"net/http"
)

// Auth errors
var (
	ErrCreateVerifyToken = New("create-verify-token", "failed to create 2fa verify token", http.StatusInternalServerError)
	ErrUpdateVerifyToken = New("update-verify-token", "failed to update 2fa verify token", http.StatusInternalServerError)
	ErrInvalidPassword   = New("invalid-password", "the user provided an invalid password", http.StatusUnauthorized)
	ErrHashPassword      = New("hash-password", "failed to hash password", http.StatusInternalServerError)
	ErrSignToken         = New("sign-token", "failed to sign jwt token", http.StatusInternalServerError)
)

// User errors
var (
	ErrMissingUserData = New("missing-user-data", "data missing to login or register", http.StatusBadRequest)
	ErrBindUser        = New("bind-user", "failed to bind to user model", http.StatusInternalServerError)
	ErrCreateUser      = New("create-user", "failed to create user", http.StatusInternalServerError)
	ErrGetUser         = New("get-user", "failed to get user", http.StatusInternalServerError)
)

// Invite errors
var (
	ErrMissingInviteData = New("missing-invite-data", "data missing to create invite", http.StatusBadRequest)
	ErrBindInvite        = New("bind-invite", "failed to bind to invite model", http.StatusInternalServerError)
	ErrCreateInvite      = New("create-invite", "failed to create invite", http.StatusInternalServerError)
)

// Room errors
var (
	ErrBindRoom        = New("bind-room", "failed to bind to room model", http.StatusInternalServerError)
	ErrMissingRoomData = New("missing-room-data", "data missing to create room", http.StatusBadRequest)
	ErrCreateRoom      = New("create-room", "failed to create room", http.StatusInternalServerError)
)

// Misc
var (
	ErrCreateAvatar = New("create-avatar", "failed to create avatar", http.StatusInternalServerError)
)

// ServiceError is a custom error returned form the service layer
type ServiceError struct {
	err     error // Contains the error
	details error // Contains detailed error message
	code    int   // HTTP status code
}

// New returns a new ServiceError
func New(err, details string, code int) *ServiceError {
	e := fmt.Errorf("%v", err)
	d := fmt.Errorf("%v", details)

	return &ServiceError{
		err:     e,
		details: d,
		code:    code,
	}
}

// Error returns the ServiceError in its string representation
func (s *ServiceError) Error() string {
	return fmt.Sprintf("status %d: '%v' %v", s.code, s.err, s.details)
}

// IsInternal returns if the error is an InternalServerError
func (s *ServiceError) IsInternal() bool {
	return s.code == http.StatusInternalServerError
}

// Err returns the short error as a string
func (s *ServiceError) Err() string {
	return s.err.Error()
}

// Details returns the detailed error message as a string
func (s *ServiceError) Details() string {
	return s.details.Error()
}

// Code returns the error code
func (s *ServiceError) Code() int {
	return s.code
}
