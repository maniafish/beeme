package sort

import (
	"fmt"
	"testing"
)

type User struct {
	Name   string `schema:"name" json:"3"`
	Passwd string `schema:"password" json:"2"`
	A      string `schema:"a" json:"4"`
	C      string `json:"1"`
	b      string `schema:"b"`
}

func TestGetFieldSortedByTag(t *testing.T) {
	fields := GetFieldsByTag(1, "schema") // int
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	fields = GetFieldsByTag("string", "schema") // string
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	u := User{"kcl", "123", "a", "c", "b"}

	fields = GetFieldsByTag(u, "schema")
	fmt.Printf("fields: %v\n", fields)

	expTags := []string{"a", "b", "name", "password"}
	expValues := []string{"a", "b", "kcl", "123"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Value.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Value.String(), expValues[i])
			return
		}
	}

	fields = GetFieldsByTag(u, "json")
	fmt.Printf("fields: %v\n", fields)

	expTags = []string{"1", "2", "3", "4"}
	expValues = []string{"c", "123", "kcl", "a"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Value.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Value.String(), expValues[i])
			return
		}
	}
}

type Foo struct {
	A int     `schema:"aa"`
	B float64 `schema:"00"`
	C string  `schema:"cc,omitempty" sort:"haha"`
	D string  `schema:"1"`
	E int     `schema:"e,omitempty"`
	F int     `schema:"f,omitempty"`
}

var foo = Foo{3, 8.33, "", "d ^", 0, 9}

func TestGetFields(t *testing.T) {
	// get tagkey "sort" as default
	fields := GetFields(foo)
	if fields.Len() != 1 {
		t.Errorf("fields length expect 1, got:%v", fields.Len())
		return
	}
}

func TestEncodeURL(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.EncodeURL("&", true)
	fmt.Println(str)
	expect := "00=8.33&1=d+%5E&aa=3&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func TestEncodeURLValOnly(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.EncodeURLValOnly("", true)
	expect := "8.33d+%5E39"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}
func TestEncode(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.Encode("&")
	fmt.Println(str)
	expect := "00=8.33&1=d ^&aa=3&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func TestEncodeValOnly(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.EncodeValOnly("")
	expect := "8.33d ^39"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func TestEncodeURLKeyValSep(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.EncodeURLKeyValSep("&", ":", true)
	expect := "00:8.33&1:d+%5E&aa:3&f:9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func BenchmarkGetFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFieldsByTag(foo, "schema")
	}
}

func TestGetNestFields(t *testing.T) {
	type SubAddr struct {
		Post   string `sort:"post"`
		Street string `sort:"street,omitempty"`
		Pass   string `sort:"-"`
		Number string
	}

	type Address struct {
		Add string `sort:"add"`
		Res string `sort:"res"`
		SubAddr
	}

	type User struct {
		ID   string `sort:"id"`
		Name string `sort:"name"`
		Address
	}

	u := User{
		ID:   "1001",
		Name: "cg",
		Address: Address{
			Add: "local^&",
			Res: "-12",
			SubAddr: SubAddr{
				Post:   "comeon. ",
				Pass:   "pass",
				Number: "No.1",
			},
		},
	}

	fields := GetNestFields(u)
	expTags := []string{"add", "id", "name", "post", "res"}
	expVals := []string{"local^&", "1001", "cg", "comeon. ", "-12"}

	t.Logf("fields: %+v", fields)

	if len(fields) != len(expTags) {
		t.Errorf("fields length expect: %v, got:%v", len(expTags), len(fields))
		return
	}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", expTags[i], f.Tag)
			return
		}
		if f.Value.String() != expVals[i] {
			t.Errorf("expect: %v, got: %v", expVals[i], f.Value.String())
			return
		}
	}
}

var foo1 = Foo{3, 8.33, "%E6%B5%8B%E8%AF%95", " Da ^测试  ", 0, 9}

func TestEncodeWithFlag(t *testing.T) {
	fields := GetFieldsByTag(foo1, "schema")
	str := fields.EncodeWithFlag("&", 0)
	expect := "00=8.33&1= Da ^测试  &aa=3&cc=%E6%B5%8B%E8%AF%95&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("", 0)
	expect = "00=8.331= Da ^测试  aa=3cc=%E6%B5%8B%E8%AF%95f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", ValOnly)
	expect = "8.33& Da ^测试  &3&%E6%B5%8B%E8%AF%95&9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", WithoutEquality)
	expect = "008.33&1 Da ^测试  &aa3&cc%E6%B5%8B%E8%AF%95&f9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", TrimSpace)
	expect = "00=8.33&1=Da ^测试&aa=3&cc=%E6%B5%8B%E8%AF%95&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", URLEncode)
	expect = "00=8.33&1=+Da+%5E%E6%B5%8B%E8%AF%95++&aa=3&cc=%25E6%25B5%258B%25E8%25AF%2595&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", URLDecode)
	expect = "00=8.33&1= Da ^测试  &aa=3&cc=测试&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", DoubleQuotation)
	expect = `00="8.33"&1=" Da ^测试  "&aa="3"&cc="%E6%B5%8B%E8%AF%95"&f="9"`
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

	str = fields.EncodeWithFlag("&", URLDecode|TrimSpace)
	expect = "00=8.33&1=Da ^测试&aa=3&cc=测试&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}

}
