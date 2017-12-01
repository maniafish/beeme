package sort

import (
	"net/url"
	"testing"
)

func TestEncodeValues(t *testing.T) {
	values := url.Values{}
	values.Set("k2", "v2")
	values.Set("k1", "")
	values.Set("k3", "")
	values.Set("k4", "v4")

	type ValueCase struct {
		Sep       string
		Omitempty bool
		Escape    map[string]bool
		Expect    string
	}
	cases := []ValueCase{
		ValueCase{"&", true, nil, "k2=v2&k4=v4"},
		ValueCase{"@", true, nil, "k2=v2@k4=v4"},
		ValueCase{"@", true, map[string]bool{"k1": true}, "k2=v2@k4=v4"},
		ValueCase{"", true, map[string]bool{"k4": true}, "k2=v2"},

		ValueCase{"&", false, nil, "k1=&k2=v2&k3=&k4=v4"},
		ValueCase{"@", false, nil, "k1=@k2=v2@k3=@k4=v4"},
		ValueCase{"@", false, map[string]bool{"k1": true}, "k2=v2@k3=@k4=v4"},
		ValueCase{"", false, map[string]bool{"k4": true}, "k1=k2=v2k3="},
	}

	for _, c := range cases {
		res := EncodeForm(values, c.Sep, c.Omitempty, c.Escape)
		if res != c.Expect {
			t.Errorf("expect: %v, got: %v", c.Expect, res)
		}
	}
}

func TestEncodeFormWithFlag(t *testing.T) {
	values := url.Values{}
	values.Set("aa", "3")
	values.Set("bb", "8.33")
	values.Set("cc", "")
	values.Set("dd", "%E6%B5%8B%E8%AF%95")
	values.Set("ee", " EE ")

	type ValueCase struct {
		Sep    string
		Flag   Flag
		Escape map[string]bool
		Expect string
	}
	cases := []ValueCase{
		ValueCase{"&", 0, nil, "aa=3&bb=8.33&cc=&dd=%E6%B5%8B%E8%AF%95&ee= EE "},
		ValueCase{"@", 0, map[string]bool{"dd": true}, "aa=3@bb=8.33@cc=@ee= EE "},
		ValueCase{"@", ValOnly, nil, "3@8.33@@%E6%B5%8B%E8%AF%95@ EE "},
		ValueCase{"@", WithoutEquality, nil, "aa3@bb8.33@cc@dd%E6%B5%8B%E8%AF%95@ee EE "},
		ValueCase{"@", TrimSpace, nil, "aa=3@bb=8.33@cc=@dd=%E6%B5%8B%E8%AF%95@ee=EE"},
		ValueCase{"@", URLEncode, nil, "aa=3@bb=8.33@cc=@dd=%25E6%25B5%258B%25E8%25AF%2595@ee=+EE+"},
		ValueCase{"@", URLDecode, nil, "aa=3@bb=8.33@cc=@dd=测试@ee= EE "},
		ValueCase{"@", DoubleQuotation, nil, `aa="3"@bb="8.33"@cc=""@dd="%E6%B5%8B%E8%AF%95"@ee=" EE "`},
		ValueCase{"@", URLDecode | TrimSpace, nil, "aa=3@bb=8.33@cc=@dd=测试@ee=EE"},
	}

	for _, c := range cases {
		res := EncodeFormWithFlag(values, c.Sep, c.Flag, c.Escape)
		if res != c.Expect {
			t.Errorf("expect: %v, got: %v", c.Expect, res)
		}
	}
}

func TestEncodeFormValOnly(t *testing.T) {
	values := url.Values{}
	values.Set("k2", "v2")
	values.Set("k1", "")
	values.Set("k3", "")
	values.Set("k4", "v4")

	type ValueCase struct {
		Sep       string
		Omitempty bool
		Escape    map[string]bool
		Expect    string
	}
	cases := []ValueCase{
		ValueCase{"&", true, nil, "v2&v4"},
		ValueCase{"@", true, nil, "v2@v4"},
		ValueCase{"@", true, map[string]bool{"k1": true}, "v2@v4"},
		ValueCase{"", true, map[string]bool{"k4": true}, "v2"},

		ValueCase{"&", false, nil, "&v2&&v4"},
		ValueCase{"@", false, nil, "@v2@@v4"},
		ValueCase{"@", false, map[string]bool{"k1": true}, "v2@@v4"},
		ValueCase{"", false, map[string]bool{"k4": true}, "v2"},
	}

	for _, c := range cases {
		res := EncodeFormValOnly(values, c.Sep, c.Omitempty, c.Escape)
		if res != c.Expect {
			t.Errorf("expect: %v, got: %v", c.Expect, res)
		}
	}
}
