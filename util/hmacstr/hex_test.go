package hmacstr

import (
	"crypto"
	"crypto/hmac"
	"fmt"
	"testing"

	_ "crypto/sha256"
)

func BenchmarkHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HexString([]byte("hello world, this is from gopher!"), "test", crypto.SHA256)
	}

}

func BenchmarkHexString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hexString2([]byte("hello world, this is from gopher!"), "test", crypto.SHA256)
	}

}
func hexString2(src []byte, key string, hashAlgo crypto.Hash) string {
	keyBytes := []byte(key)
	mac := hmac.New(hashAlgo.New, keyBytes)
	mac.Write(src)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func TestHexString(t *testing.T) {

	tests := []struct {
		msg  []byte
		key  string
		hash crypto.Hash
		sum  string
	}{
		{
			msg:  []byte("hello world"),
			key:  "test",
			hash: crypto.SHA256,
			sum:  "d1596e0d4280f2bd2d311ce0819f23bde0dc834d8254b92924088de94c38d922",
		},
	}

	for _, test := range tests {
		if HexString(test.msg, test.key, test.hash) != hexString2(test.msg, test.key, test.hash) {
			t.Errorf("two function result not equal")
		}
		sum := HexString(test.msg, test.key, test.hash)
		if sum != test.sum {
			t.Errorf("result %s != sum %s", sum, test.sum)
		}

	}

}
