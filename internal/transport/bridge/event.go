// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import "github.com/pion/webrtc/v3"

type Event struct {
	Type string `json:"string"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
	User      *PublicUser                `json:"user,omitempty"`
}

const (
	EventOffer     string = "offer"
	EventAnswer    string = "answer"
	EventCandidate string = "candidate"
	EventMute      string = "mute"
	EventUnmute    string = "unmute"
	EventJoin      string = "user-join"
	EventLeave     string = "user-leave"
)
