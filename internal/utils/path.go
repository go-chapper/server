package utils

import (
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// Abs returns the absolute path of p
func Abs(p string) (string, error) {
	return filepath.Abs(p)
}

// Join joins the paths
func Join(p ...string) string {
	return filepath.Join(p...)
}

func ConfigPath() (string, error) {
	if runtime.GOOS == "linux" {
		return filepath.Abs("/etc/chapper")
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".chapper"), nil
}

func UserHome() (string, error) {
	return homedir.Dir()
}
