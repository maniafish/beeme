package md5str

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// HexString return hex hash string
func HexString(text string) string {
	sum := md5.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

// Base64String hash text with md5, return base64 string
func Base64String(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
