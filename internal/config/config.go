// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package config provides the top-level config struct
package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"

	"chapper.dev/server/internal/utils"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log     LogOptions
	Turn    TurnOptions
	Store   StoreOptions
	Router  RouterOptions
	General GeneralOptions
}

type LogOptions struct {
	Path   string `toml:"PATH"`
	Prefix string `toml:"PREFIX"`
}

type TurnOptions struct {
	Domain   string `toml:"DOMAIN"`
	PublicIP string `toml:"PUBLIC_IP"`
	Port     int    `toml:"PORT"`
}

type StoreOptions struct {
	User     string `toml:"USER"`
	Password string `toml:"PASSWORD"`
	Database string `toml:"DATABASE"`
	Host     string `toml:"HOST"`
	Port     int    `toml:"PORT"`
}

type RouterOptions struct {
	Port          int    `toml:"PORT"`
	Domain        string `toml:"DOMAIN"`
	WebPath       string `toml:"WEB_PATH"`
	AvatarPath    string `toml:"AVATAR_PATH"`
	JWTSigningKey string `toml:"JWT_SIGNING_KEY"`
	OTPIssuer     string `toml:"OTP_ISSUER"`
	EnableGZIP    bool   `toml:"ENABLE_GZIP"`
}

type GeneralOptions struct {
	Name           string `toml:"NAME"`
	EnableRegister bool   `toml:"ENABLE_REGISTER"`
	DisableBanner  bool   `toml:"DISABLE_BANNER"`
}

// New returns a new config struct
func New() *Config {
	return &Config{}
}

// NewDefault returns a default config
func NewDefault() *Config {
	if runtime.GOOS == "linux" {
		return &Config{
			Log: LogOptions{
				Path:   "/var/log/chapper/chapper.log",
				Prefix: "Chapper",
			},
			Store: StoreOptions{
				User:     "root",
				Password: "",
				Database: "chapper",
				Host:     "127.0.0.1",
				Port:     3306,
			},
			Router: RouterOptions{
				Port:          8080,
				Domain:        "",
				WebPath:       "/var/www/chapper/app",
				AvatarPath:    "/var/www/chapper/avatar",
				JWTSigningKey: "",
				OTPIssuer:     "Chapper",
				EnableGZIP:    true,
			},
			General: GeneralOptions{
				Name:           "Chapper",
				EnableRegister: true,
				DisableBanner:  false,
			},
		}
	}

	return &Config{
		Log: LogOptions{
			Path:   "",
			Prefix: "Chapper",
		},
		Store: StoreOptions{
			User:     "",
			Password: "",
			Database: "",
			Host:     "",
			Port:     3306,
		},
		Router: RouterOptions{
			Port:          8080,
			Domain:        "",
			WebPath:       "",
			AvatarPath:    "",
			JWTSigningKey: "",
			OTPIssuer:     "Chapper",
			EnableGZIP:    true,
		},
		General: GeneralOptions{
			Name:           "Chapper",
			EnableRegister: true,
			DisableBanner:  false,
		},
	}
}

// Read reads a .toml config file from path and returns a config struct or an error
func (c *Config) Read(path string) error {
	path, err := utils.Abs(path)
	if err != nil {
		return err
	}

	_, err = toml.DecodeFile(path, c)
	if err != nil {
		return err
	}

	// Validate the config
	return c.Validate()
}

func (c *Config) Write(path string) error {
	path, err := utils.Abs(path)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	err = toml.NewEncoder(buffer).Encode(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0644)
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.Router.JWTSigningKey == "" {
		// TODO: Use logging
		fmt.Println("WARNING [Config] JWT_SIGNING_KEY cannot be empty. Fallback to generated one")
		key, err := utils.RandomCryptoString(32)
		if err != nil {
			return err
		}
		c.Router.JWTSigningKey = key
	}

	if c.Router.OTPIssuer == "" {
		// Fallback to default issuer
		fmt.Println("WARNING [Config] OTP_ISSUER cannot be empty. Fallback to 'Chapper'")
		c.Router.OTPIssuer = "Chapper"
	}

	if c.Log.Prefix == "" {
		c.Log.Prefix = "Chapper"
	}

	return nil
}
