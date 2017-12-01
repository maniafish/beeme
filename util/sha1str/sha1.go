package sha1str

import (
	"crypto/sha1"
	"encoding/hex"
)

// HexString return hex hash string
func HexString(text string) string {
	sum := sha1.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}
