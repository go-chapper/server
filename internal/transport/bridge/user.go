// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v2"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 51200
)

var (
	ErrNoSuchEvent   = errors.New("no such event")
	ErrChanClosed    = errors.New("channel closed")
	ErrInvalidTrack  = errors.New("track is nil")
	ErrInvalidPacket = errors.New("packet is nil")
)

// User keeps track of a websocket connection for signaling, a WebRTC peer connection and
// connects hub, room and user
type User struct {
	ID            string                   // A unique ID of the user
	room          *Room                    // The room the user is in
	conn          *websocket.Conn          // The underlying websocket connection to exchange data (signaling)
	send          chan []byte              // Channel for outbound messages
	pc            *webrtc.PeerConnection   // WebRTC peer connection
	inTracks      map[uint32]*webrtc.Track // Incoming tracks (microphone)
	inTracksLock  sync.RWMutex             // Incoming tracks lock
	outTracks     map[uint32]*webrtc.Track // Rest of the room's tracks
	outTracksLock sync.RWMutex

	rtpCh chan *rtp.Packet

	stop bool
	info UserInfo
}

// PublicUser is the public representation of a user. It only includes necessary data
type PublicUser struct {
	ID string `json:"id"`
	UserInfo
}

// UserInfo holds some user information
type UserInfo struct {
	Username string `json:"username"`
	Mute     bool   `josn:"mute"`
}

// NewUser creates and returns a new user
func NewUser(username string, conn *websocket.Conn, pc *webrtc.PeerConnection, room *Room) *User {
	return &User{
		ID:        uuid.New().String(),
		room:      room,
		conn:      conn,
		send:      make(chan []byte, 256),
		pc:        pc,
		inTracks:  make(map[uint32]*webrtc.Track),
		outTracks: make(map[uint32]*webrtc.Track),
		rtpCh:     make(chan *rtp.Packet, 100),
		info: UserInfo{
			Username: username,
			Mute:     false,
		},
	}
}

// ToPublic returns the public representation of a user
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:       u.ID,
		UserInfo: u.info,
	}
}

func (u *User) AddListeners() error {
	u.pc.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		if iceCandidate != nil {
			err := u.sendCandidate(iceCandidate)
			if err != nil {
				// Log
				return
			}
		}
	})

	return nil
}

// startRead starts continuously reading from the user's websocket
func (u *User) startRead() {
	defer func() {
		u.stop = true
		u.pc.Close()
		u.room.Leave(u)
		u.conn.Close()
	}()

	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		event := Event{}
		err := u.conn.ReadJSON(&event)
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				log.Println(err)
			}
			break
		}

		go func() {
			err := u.handleEvent(event)
			if err != nil {
				log.Println(err)
				// u.SendErr(err)
			}
		}()
	}
}

// startWrite starts to continuously listen for incoming messages to be written to the
// websocket connection
func (u *User) startWrite() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		u.stop = true
		u.conn.Close()
	}()

	for {
		select {
		case message, ok := <-u.send:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The bridge closed the channel.
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sendEvent sends JSON to the user's websocket
func (u *User) sendEvent(event Event) error {
	json, err := json.Marshal(event)
	if err != nil {
		return err
	}
	u.send <- json
	return nil
}

// handleEvent handles the different signaling / WebRTC events
func (u *User) handleEvent(event Event) error {
	switch event.Type {
	case EventOffer:
		if event.Offer == nil {
			// Handle error
		}

		return u.handleOffer(*event.Offer)
	case EventAnswer:
		if event.Answer == nil {
			// Handle error
		}

		return u.pc.SetRemoteDescription(*event.Answer)
	case EventCandidate:
		if event.Candidate == nil {
			// Handle error
		}

		return u.pc.AddICECandidate(*event.Candidate)
	case EventMute:
		u.info.Mute = true
		return u.room.BroadcastEventMute(u)
	case EventUnmute:
		u.info.Mute = false
		return u.room.BroadcastEventUnmute(u)
	default:
		return ErrNoSuchEvent
	}
}

// handleOffer handles the incoming session description / offer
func (u *User) handleOffer(offer webrtc.SessionDescription) error {
	// Add receive only transciever. This is the users microphone
	if len(u.pc.GetTransceivers()) == 0 {
		_, err := u.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RtpTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		})
		if err != nil {
			return err
		}
	}

	// Set remote session description
	err := u.pc.SetRemoteDescription(offer)
	if err != nil {
		return err
	}

	return u.sendAnswer()
}

// sendCandidate sends ice candidate to peer
func (u *User) sendCandidate(iceCandidate *webrtc.ICECandidate) error {
	if iceCandidate == nil {
		return errors.New("nil ice candidate")
	}
	iceCandidateInit := iceCandidate.ToJSON()
	err := u.sendEvent(Event{Type: "candidate", Candidate: &iceCandidateInit})
	if err != nil {
		return err
	}
	return nil
}

// answer creates a session description answer
func (u *User) answer() (webrtc.SessionDescription, error) {
	answer, err := u.pc.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = u.pc.SetLocalDescription(answer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	return answer, nil
}

// sendAnswer creates answer and sends it via websocket
func (u *User) sendAnswer() error {
	answer, err := u.answer()
	if err != nil {
		return err
	}
	return u.sendEvent(Event{Type: "answer", Answer: &answer})
}

// receiveInTrackRTP receive all incoming tracks' rtp and sent to one channel
func (u *User) receiveInTrackRTP(remoteTrack *webrtc.Track) {
	for {
		if u.stop {
			return
		}
		rtp, err := remoteTrack.ReadRTP()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("rtp err => %v", err)
		}
		u.rtpCh <- rtp
	}
}

// ReadRTP reads rtp packets
func (u *User) ReadRTP() (*rtp.Packet, error) {
	rtp, ok := <-u.rtpCh
	if !ok {
		return nil, ErrChanClosed
	}
	return rtp, nil
}

// WriteRTP sends rtp packets to user's outgoing tracks
func (u *User) WriteRTP(pkt *rtp.Packet) error {
	if pkt == nil {
		return ErrInvalidPacket
	}
	u.outTracksLock.RLock()
	track := u.outTracks[pkt.SSRC]
	u.outTracksLock.RUnlock()

	if track == nil {
		log.Printf("WebRTCTransport.WriteRTP track==nil pkt.SSRC=%d", pkt.SSRC)
		return ErrInvalidTrack
	}

	// log.Debugf("WebRTCTransport.WriteRTP pkt=%v", pkt)
	err := track.WriteRTP(pkt)
	if err != nil {
		// log.Errorf(err.Error())
		// u.writeErrCnt++
		return err
	}
	return nil
}

// AddTrack adds track to peer connection
func (u *User) AddTrack(ssrc uint32) error {
	track, err := u.pc.NewTrack(webrtc.DefaultPayloadTypeOpus, ssrc, string(ssrc), string(ssrc))
	if err != nil {
		return err
	}
	if _, err := u.pc.AddTrack(track); err != nil {
		log.Println("ERROR Add remote track as peerConnection local track", err)
		return err
	}

	u.outTracksLock.Lock()
	u.outTracks[track.SSRC()] = track
	u.outTracksLock.Unlock()
	return nil
}
