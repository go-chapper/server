package logger

import (
	"log"
	"os"
	"strings"

	"git.web-warrior.de/go-chapper/server/internal/config"
	"git.web-warrior.de/go-chapper/server/internal/utils"
)

// New sets up a new std logger which writes to file at 'path' with 'prefix'
func New(c config.LogOptions) (*os.File, error) {
	path, err := utils.Abs(c.Path)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return f, err
	}

	prefix := strings.TrimSpace(c.Prefix) + " | "

	log.SetOutput(f)
	log.SetPrefix(prefix)
	return f, nil
}
