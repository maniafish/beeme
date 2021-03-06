package md5str

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

func BenchmarkHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HexString("hello world, this is from gopher!")
	}

}

func BenchmarkHexString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hexString2("hello world, this is from gopher!")
	}

}
func BenchmarkHexString3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hexString3("hello world, this is from gopher!")
	}

}

func hexString2(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hexString3(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func TestHexString(t *testing.T) {

	tests := []struct {
		msg string
		sum string
	}{
		{
			msg: "hello world",
			sum: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
	}

	for _, test := range tests {
		if HexString(test.msg) != hexString2(test.msg) {
			t.Errorf("two function result not equal")
		}

		sum := HexString(test.msg)
		if sum != test.sum {
			t.Errorf("result %s != sum %s", sum, test.sum)
		}

	}

}
