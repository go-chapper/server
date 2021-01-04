// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package constants

// AvailabilityState defines the availability state a user can be in
type AvailabilityState string

const (
	Online    AvailabilityState = "online"
	Busy      AvailabilityState = "busy"
	Away      AvailabilityState = "away"
	Offline   AvailabilityState = "offline"
	Invisible AvailabilityState = "invisible"
)
