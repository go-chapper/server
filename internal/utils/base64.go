// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"encoding/base64"
	"encoding/json"
)

// DecodeBase64 decodes a base64 'input' and unmarshals it into 'target'
func DecodeBase64(input string, target interface{}) error {
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, target)
}

// EncodeBase64 encodes 'input' to base64 string
func EncodeBase64(input interface{}) (string, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
