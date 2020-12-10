// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import "github.com/pion/webrtc/v2"

// Event represents a WebRTC signaling message / event
type Event struct {
	Type string `json:"type"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
	User      *PublicUser                `json:"user,omitempty"`
}

const (
	// TypeOffer describes the event type 'offer'
	TypeOffer string = "offer"

	// TypeAnswer describes the event type 'answer'
	TypeAnswer string = "answer"

	// TypeCandidate describes the event type 'candidate'
	TypeCandidate string = "candidate"

	// TypeMute describes the event type 'mute'
	TypeMute string = "mute"

	// TypeUnmute describes the event type 'unmute'
	TypeUnmute string = "unmute"

	// TypeUser describes the event type 'user'
	TypeUser string = "user"

	// TypeUserJoin describes the event type 'user-join'
	TypeUserJoin string = "user-join"

	// TypeUserLeave describes the event type 'user-leave'
	TypeUserLeave string = "user-leave"
)
