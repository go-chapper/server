// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package schemas

const Rooms = `
CREATE TABLE IF NOT EXISTS rooms (
	hash VARCHAR(32) NOT NULL,
	name VARCHAR(100) NOT NULL,
	type VARCHAR(10) DEFAULT NULL,
	description TEXT DEFAULT NULL,
	PRIMARY KEY (hash)
);
`
