package hmacstr

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
)

// Base64String return hmac result in base64 encoding string
func Base64String(src []byte, key string, hashAlgo crypto.Hash) string {
	keyBytes := []byte(key)
	mac := hmac.New(hashAlgo.New, keyBytes)
	mac.Write(src)
	signResult := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signResult
}
