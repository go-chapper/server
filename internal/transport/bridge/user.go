// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type User struct {
	ID            string                         // A unique ID of the user
	room          *Room                          // The room the user is in
	conn          *websocket.Conn                // The underlying websocket connection to exchange data (signaling)
	send          chan []byte                    // Channel for outbound messages
	pc            *webrtc.PeerConnection         // WebRTC peer connection
	inTracks      map[uint32]*webrtc.TrackRemote // Incoming tracks (microphone)
	inTracksLock  sync.RWMutex                   // Incoming tracks lock
	outTracks     map[uint32]*webrtc.TrackRemote // Rest of the room's tracks
	outTracksLock sync.RWMutex

	rtpCh chan *rtp.Packet

	stop bool
	info UserInfo
}

type PublicUser struct {
	ID string `json:"id"`
	UserInfo
}

type UserInfo struct {
	Username string `json:"username"`
	Mute     bool   `josn:"mute"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:       u.ID,
		UserInfo: u.info,
	}
}
