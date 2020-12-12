// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// ToNullTime returns t as nullable time
func ToNullTime(t time.Time) null.Time {
	return null.TimeFrom(t)
}
