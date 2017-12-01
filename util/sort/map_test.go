package sort

import (
	"testing"
)

func genMap() (map[string]interface{}, string, string, string, string) {
	m := make(map[string]interface{})
	m["a"] = 23
	m["b"] = 23.2100
	m["c"] = 1234567890
	m["0"] = nil
	m["1"] = "string *^"
	m["2"] = true

	ret := "0=&1=string *^&2=true&a=23&b=23.21&c=1234567890"
	retURL := "0=&1=string+%2A%5E&2=true&a=23&b=23.21&c=1234567890"
	retValOnly := "&string *^&true&23&23.21&1234567890"
	retURLValOnly := "&string+%2A%5E&true&23&23.21&1234567890"
	return m, ret, retURL, retValOnly, retURLValOnly
}

func getStringMap() (map[string]string, string, string) {
	m := make(map[string]string)
	m["a"] = "23"
	m["b"] = "23.2100"
	m["c"] = "1234567890"
	m["1"] = "string *^"
	m["2"] = "true"
	m["0"] = ""

	ret := "0=&1=string *^&2=true&a=23&b=23.2100&c=1234567890"
	retURL := "0=&1=string+%2A%5E&2=true&a=23&b=23.2100&c=1234567890"
	return m, ret, retURL
}

func TestMap(t *testing.T) {
	m, ret, retURL, retValOnly, retURLValOnly := genMap()

	ret1 := EncodeMap(m, "&", false)
	retURL2 := EncodeMap(m, "&", true)
	retValOnly2 := EncodeMapValOnly(m, "&", false)
	retURLValOnly2 := EncodeMapValOnly(m, "&", true)

	if ret != ret1 {
		t.Errorf("expect: %v, got: %v", ret, ret1)
		return
	}
	if retURL != retURL2 {
		t.Errorf("expect: %v, got: %v", retURL, retURL2)
		return
	}
	if retValOnly != retValOnly2 {
		t.Errorf("expect: %v, got: %v", retValOnly, retValOnly2)
		return
	}
	if retURLValOnly != retURLValOnly2 {
		t.Errorf("expect: %v, got: %v", retURLValOnly, retURLValOnly2)
		return
	}
}

func TestStringMap(t *testing.T) {
	m, ret, retURL := getStringMap()

	ret1 := EncodeMapString(m, "&", false)
	retURL2 := EncodeMapString(m, "&", true)

	if ret != ret1 {
		t.Errorf("expect: %v, got: %v", ret, ret1)
		return
	}
	if retURL != retURL2 {
		t.Errorf("expect: %v, got: %v", retURL, retURL2)
		return
	}
}

func BenchmarkFormatStr(b *testing.B) {
	m, _, _, _, _ := genMap()
	for i := 0; i < b.N; i++ {
		for _, v := range m {
			formatStr(v)
		}
	}
}

func BenchmarkEncodeMap(b *testing.B) {
	m, _, _, _, _ := genMap()
	for i := 0; i < b.N; i++ {
		EncodeMap(m, "&", false)
		EncodeMap(m, "&", true)
		EncodeMapValOnly(m, "&", false)
		EncodeMapValOnly(m, "&", true)
	}
}
