// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package schemas

const Servers = `
CREATE TABLE IF NOT EXISTS servers (
	hash VARCHAR(32) NOT NULL,
	name VARCHAR(100) NOT NULL,
	description TEXT DEFAULT NULL,
	image VARCHAR(32) DEFAULT NULL,
	PRIMARY KEY (hash)
);
`
