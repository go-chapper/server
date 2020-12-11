// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package services provides different services which combine functions from modules and
// the store interface
package services

import (
	"fmt"
	"net/http"
)

// ServiceError is a custom error returned form the service layer
type ServiceError struct {
	err     error // Contains the error
	details error // Contains detailed error message
	code    int   // HTTP status code
}

// NewError returns a new ServiceError
func NewError(err, details string, code int) *ServiceError {
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
	return fmt.Sprintf("Status %d: '%v' %v", s.code, s.err, s.details)
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
