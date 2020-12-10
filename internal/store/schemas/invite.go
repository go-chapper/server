// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package schemas

const Invites = `
CREATE TABLE IF NOT EXISTS invites (
	hash VARCHAR(32) NOT NULL,
	creayted_by VARCHAR(100) NOT NULL,
	server VARCHAR(100) NOT NULL,
	one_time_use BOOLEAN DEFAULT false,
	expires_at DATETIME DEFAULT NULL,
	PRIMARY KEY (hash)
);
`
