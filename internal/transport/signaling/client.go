// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities to handle WebRTC signaling
package signaling

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

const (
	writeWait = 2 * time.Second
)

// Client is an active client connection to the signaling service
type Client struct {
	conn  *websocket.Conn
	hub   *Hub
	info  ClientInfo
	send  chan Message
	close chan string
}

// ClientInfo holds information about the connected client
type ClientInfo struct {
	ID   xid.ID
	User string
}

// NewClient creates a new client and returns it
func NewClient(conn *websocket.Conn, hub *Hub, user string) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
		info: ClientInfo{
			ID:   xid.New(),
			User: user,
		},
		send: hub.broadcast,
	}
}

func (c *Client) Close() {
	c.conn.Close()
	go func() {
		// c.send <- ClientMessage{
		// 	Info: c.info,
		// 	// Incoming: &Disconnected{},
		// }
	}()
}

func (c *Client) listenRead(pongInterval time.Duration) {
	defer c.Close()
	_ = c.conn.SetReadDeadline(time.Now().Add(pongInterval))
	c.conn.SetPongHandler(func(appData string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	for {
		t, m, err := c.conn.NextReader()
		if err != nil {
			// c.printWebSocketError("read", err)
			return
		}
		if t == websocket.BinaryMessage {
			_ = c.conn.CloseHandler()(websocket.CloseUnsupportedData, fmt.Sprintf("unsupported binary message type: %s", err))
			return
		}

		incoming, err := c.hub.readIncoming(m)
		if err != nil {
			_ = c.conn.CloseHandler()(websocket.CloseNormalClosure, fmt.Sprintf("malformed message: %s", err))
			return
		}

		c.send <- incoming
	}
}

func (c *Client) listenWrite(pingInterval time.Duration) {
	pingTicker := time.NewTicker(pingInterval)

	closeConnection := func() {
		c.Close()
		pingTicker.Stop()
	}

	defer closeConnection()

	for {
		select {
		case reason := <-c.send:
			fmt.Println(reason)
			closeConnection()
		case message := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			outgoing, err := c.hub.writeOutgoing(message)
			if err != nil {
				closeConnection()
				continue
			}

			if err := c.conn.WriteJSON(outgoing); err != nil {
				closeConnection()
			}
		case <-pingTicker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				closeConnection()
			}
		}
	}
}
