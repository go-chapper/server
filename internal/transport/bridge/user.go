// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import (
	"encoding/json"
	"errors"
	"fmt"
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

// AddListeners adds all neccesary WebRTC PC event handlers
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

	u.pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateConnected {
			log.Println("user joined")
			tracks := u.getRoomTracks()
			fmt.Println("attach ", len(tracks), "tracks to new user")
			u.log("new user add tracks", len(tracks))
			for _, track := range tracks {
				err := u.AddTrack(track.SSRC())
				if err != nil {
					log.Println("ERROR Add remote track as peerConnection local track", err)
					panic(err) // NOTE(Techassi): Should we really panic here?
				}
			}
			err := u.SendOffer()
			if err != nil {
				panic(err) // NOTE(Techassi): Should we really panic here?
			}
		} else if connectionState == webrtc.ICEConnectionStateDisconnected ||
			connectionState == webrtc.ICEConnectionStateFailed ||
			connectionState == webrtc.ICEConnectionStateClosed {

			u.stop = true
			senders := u.pc.GetSenders()
			for _, roomUser := range u.room.GetParticipants(u) {
				// fmt.Println("removing tracks from user")
				u.log("removing tracks from user")
				for _, sender := range senders {
					ssrc := sender.Track().SSRC()

					roomUserSenders := roomUser.pc.GetSenders()
					for _, roomUserSender := range roomUserSenders {
						if roomUserSender.Track().SSRC() == ssrc {
							err := roomUser.pc.RemoveTrack(roomUserSender)
							if err != nil {
								panic(err)
							}
						}
					}
				}
			}

		}
	})

	u.pc.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		u.log(
			"peerConnection.OnTrack",
			fmt.Sprintf("track has started, of type %d: %s, ssrc: %d \n", remoteTrack.PayloadType(), remoteTrack.Codec().Name, remoteTrack.SSRC()),
		)
		if _, alreadyAdded := u.inTracks[remoteTrack.SSRC()]; alreadyAdded {
			u.log("user.inTrack != nil", "already handled")
			return
		}

		u.inTracks[remoteTrack.SSRC()] = remoteTrack
		for _, roomUser := range u.room.GetParticipants(u) {
			log.Println("add remote track", fmt.Sprintf("(user: %s)", u.ID), "track to user ", roomUser.ID)
			if err := roomUser.AddTrack(remoteTrack.SSRC()); err != nil {
				log.Println(err)
				continue
			}
			err := roomUser.SendOffer()
			if err != nil {
				panic(err)
			}
		}
		go u.receiveInTrackRTP(remoteTrack)
		go u.broadcastIncomingRTP()
	})

	return nil
}

// GetRoomTracks returns list of room incoming tracks
func (u *User) getRoomTracks() []*webrtc.Track {
	tracks := []*webrtc.Track{}
	for _, user := range u.room.GetUsersList() {
		for _, track := range user.inTracks {
			tracks = append(tracks, track)
		}
	}
	return tracks
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

// SendEventUser sends user to client to identify himself
func (u *User) SendEventUser() error {
	return u.sendEvent(Event{Type: "user", User: u.ToPublic()})
}

// SendEventRoom sends room to client with users except me
// func (u *User) SendEventRoom() error {
// 	return u.sendEvent(Event{Type: "room", Room: u.room.Wrap(u)})
// }

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

// Offer return a offer
func (u *User) Offer() (webrtc.SessionDescription, error) {
	offer, err := u.pc.CreateOffer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	err = u.pc.SetLocalDescription(offer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	return offer, nil
}

// SendOffer creates webrtc offer
func (u *User) SendOffer() error {
	offer, err := u.Offer()
	err = u.sendEvent(Event{Type: "offer", Offer: &offer})
	if err != nil {
		panic(err)
	}
	return nil
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

func (u *User) broadcastIncomingRTP() {
	for {
		rtp, err := u.ReadRTP()
		if err != nil {
			panic(err)
		}
		for _, user := range u.room.GetParticipants(u) {
			err := user.WriteRTP(rtp)
			if err != nil {
				// panic(err)
				fmt.Println(err)
			}
		}
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

func (u *User) log(msg ...interface{}) {
	log.Println(
		fmt.Sprintf("user %s:", u.ID),
		fmt.Sprint(msg...),
	)
}

// GetOutTracks return outgoing tracks
func (u *User) GetOutTracks() map[uint32]*webrtc.Track {
	u.outTracksLock.RLock()
	defer u.outTracksLock.RUnlock()
	return u.outTracks
}

// Watch for debug
func (u *User) Watch() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		if u.stop {
			ticker.Stop()
			return
		}
		fmt.Println("ID:", u.ID, "out: ", u.GetOutTracks())
	}
}
