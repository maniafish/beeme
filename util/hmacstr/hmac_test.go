package hmacstr

import (
	"crypto"
	"testing"

	_ "crypto/sha256"
)

func TestBase64String(t *testing.T) {

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
			sum:  "0VluDUKA8r0tMRzggZ8jveDcg02CVLkpJAiN6Uw42SI=",
		},
	}

	for _, test := range tests {
		sum := Base64String(test.msg, test.key, test.hash)
		if sum != test.sum {
			t.Errorf("result %s != sum %s", sum, test.sum)
		}

	}

}
