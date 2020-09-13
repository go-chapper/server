// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package turn provides utilities to run a TURN server
package turn

import (
	"errors"
	"net"
	"strconv"

	"github.com/pion/turn/v2"
)

type TURN struct {
	PublicIP string
	Realm    string
	Port     int
	Protocol string
	server   *turn.Server
}

var (
	ErrPublicIPEmpty       = errors.New("Public IP empty")
	ErrUnsupportedProtocol = errors.New("Protocol unsupported")
)

// New returns a new TURN instance or an error
// 'publicIP' is the public IP which can be contacted by users via the internet
// 'realm' is the realm this TURN server runs under, usually your domain, defaults to
// chapper.dev
// 'port' the port this TURN server listens on, defaults to 50554
func New(publicIP, realm, protocol string, port int) (*TURN, error) {
	if publicIP == "" {
		return nil, ErrPublicIPEmpty
	}

	if protocol != "udp4" && protocol != "tcp4" {
		return nil, ErrUnsupportedProtocol
	}

	if port == 0 {
		port = 50554
	}

	if realm == "" {
		realm = "chapper.dev"
	}

	return &TURN{
		PublicIP: publicIP,
		Realm:    realm,
		Port:     port,
		Protocol: protocol,
	}, nil
}

// Run runs the TURN server
func (t *TURN) Run() error {
	listener, err := net.ListenPacket(t.Protocol, "0.0.0.0:"+strconv.Itoa(t.Port))
	if err != nil {
		return err
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: t.Realm,
		// Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		// AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
		// 	if key, ok := usersMap[username]; ok {
		// 		return key, true
		// 	}
		// 	return nil, false
		// },
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: listener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(t.PublicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",               // But actually be listening on every interface
				},
			},
		},
	})
	if err != nil {
		return err
	}

	t.server = s
	return nil
}

// Close closes the TURN server
func (t *TURN) Close() error {
	return t.server.Close()
}
