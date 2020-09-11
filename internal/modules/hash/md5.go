package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 returns the MD5 hash of 'payload'
func MD5(payload string) string {
	h := md5.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
