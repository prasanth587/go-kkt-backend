package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// SHAEncoding encodes the string
func SHAEncoding(target string) (output string) {
	buf := []byte(target)

	encrypted := sha256.New()
	encrypted.Write(buf)

	output = base64.StdEncoding.EncodeToString(encrypted.Sum(nil))

	return
}
