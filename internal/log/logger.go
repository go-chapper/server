// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package log provides logging utilities
package log

import (
	"log"
	"os"
	"strings"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/utils"
)

// Logger wraps a log file and the standard logger
type Logger struct {
	file   *os.File
	logger *log.Logger
}

// Context describes a log context
type Context string

// New sets up a new logger which writes to file at 'path' with 'prefix'
func New(c config.LogOptions) (*Logger, error) {
	path, err := utils.Abs(c.Path)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	prefix := strings.TrimSpace(c.Prefix) + " | "
	logger := log.New(f, prefix, log.LstdFlags)

	return &Logger{
		file:   f,
		logger: logger,
	}, nil
}

// NewContext returns a new context
func NewContext(ctx string) Context {
	return Context(ctx)
}

// Fatal logs an error via the underlying standard logger and exists with code 1
func (l *Logger) Fatal(err error) {
	l.Error(err)
	os.Exit(1)
}

// Error logs an error via the underlying standard logger
func (l *Logger) Error(err error) {
	l.logger.Printf("[E] %v", err)
}

// Errorc logs an error with context via the underlying standard logger
func (l *Logger) Errorc(ctx Context, err error) {
	l.logger.Printf("[E] [%s] %v", ctx, err)
}

// CErrorf logs an error with context via the underlying standard logger
// func (l *Logger) CErrorf(ctx Context, format string, v ...interface{}) {
// 	s := fmt.Sprintf("[E] %s: ", ctx)
// 	l.logger.Printf(s+format, v...)
// }

// Errorf logs an error via the underlying standard logger
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("[E] "+format, v...)
}

// Info logs a info message via the underlying standard logger
func (l *Logger) Info(msg string) {
	l.logger.Printf("[I] %s", msg)
}

// Infoc logs a info message with context via the underlying standard logger
func (l *Logger) Infoc(ctx Context, msg string) {
	l.logger.Printf("[I] [%s] %s", ctx, msg)
}

// Infof logs a info message via the underlying standard logger
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Printf("[I] "+format, v...)
}
