package hmacstr

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
)

// HexString return hmac result in hex encoding
func HexString(src []byte, key string, hashAlgo crypto.Hash) string {
	keyBytes := []byte(key)
	mac := hmac.New(hashAlgo.New, keyBytes)
	mac.Write(src)
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum[:])
}
