// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package schemas

const Users = `
CREATE TABLE IF NOT EXISTS users (
	username VARCHAR(100) NOT NULL,
	password VARCHAR(512) NOT NULL,
	email VARCHAR(100) DEFAULT NULL,
	publickey VARCHAR(1000) NOT NULL,
	twofa_secret VARCHAR(16) DEFAULT NULL,
	twofa_verify VARCHAR(16) DEFAULT NULL,
	PRIMARY KEY (username)
);
`
