package sort

// package sort provides primitives for sorting method in gas system

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

// TagField represents field value along with tag in any struct
type TagField struct {
	Tag   string
	Value reflect.Value
}

// TagFields wraps TagField slice, and provides Encode method.
// TagFields can be sort by tag in alphabet order.
type TagFields []TagField

// Flag type for EncodeWithFlag
type Flag int

// Flags
const (
	// ValOnly 只拼接参数值
	ValOnly = 1 << iota
	// WithoutEquality 不拼接"="，即key1value1&key2value2
	WithoutEquality
	// TrimSpace 去除参数值首尾空格
	TrimSpace
	// URLEncode 对参数值做urlencode
	URLEncode
	// URLDecode 对参数值做urldecode
	URLDecode
	// DoubleQuotation 对参数值加上双引号
	DoubleQuotation
	// OmitEmpty 只拼接非空的参数值对
	OmitEmpty
)

// Len for sort
func (t TagFields) Len() int {
	return len(t)
}

// Swap for sort
func (t TagFields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less for sort
func (t TagFields) Less(i, j int) bool {
	return strings.Compare(t[i].Tag, t[j].Tag) < 0
}

// Encode encodes TagFields into string by joining [tag=value, tag=value, ...]
// seperated by given sep.
//
// Example:
//		given sep=&, t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
//		"tag1=val1&tag2=val2&...&tagn=valn" is return
func (t TagFields) Encode(sep string) string {
	return t.EncodeURL(sep, false)
}

// EncodeValOnly encodes TagFields into string by joining [value, value, ...]
// seperated by given sep
//
// Example:
//		given sep=&, t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
//		"val1&val2&...&valn" is return
func (t TagFields) EncodeValOnly(sep string) string {
	return t.EncodeURLValOnly(sep, false)
}

// EncodeKeyValSep encodes TagFields into string by joining [value, value, ...]
// seperated by given sep
//
// Example:
//		given sep1="&", sep2=":" t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
//		"tag1:val1&tag2:val2&...&tagn:valn" is return
func (t TagFields) EncodeKeyValSep(sep1, sep2 string) string {
	return t.EncodeURLKeyValSep(sep1, sep2, false)
}

//EncodeURLKeyValSep is similar to EncodeKeyValSep, isURLEncode means whether the value need to be urlencode
func (t TagFields) EncodeURLKeyValSep(sep1, sep2 string, isURLEncode bool) string {
	var buffer bytes.Buffer
	for i, f := range t {
		if i != 0 {
			buffer.WriteString(sep1)
		}
		buffer.WriteString(f.Tag)
		buffer.WriteString(sep2)
		if isURLEncode {
			buffer.WriteString(url.QueryEscape(fmt.Sprintf("%v", f.Value)))
		} else {
			buffer.WriteString(fmt.Sprintf("%v", f.Value))
		}
	}
	return buffer.String()
}

//EncodeURL is similar to Encode, isURLEncode means whether the value need to be urlencode
func (t TagFields) EncodeURL(sep string, isURLEncode bool) string {
	var buffer bytes.Buffer
	for i, f := range t {
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(f.Tag)
		buffer.WriteString("=")
		if isURLEncode {
			buffer.WriteString(url.QueryEscape(fmt.Sprintf("%v", f.Value)))
		} else {
			buffer.WriteString(fmt.Sprintf("%v", f.Value))
		}
	}
	return buffer.String()
}

//EncodeURLValOnly is similar to EncodeURL, isURLEncode means whether the value need to be urlencode
func (t TagFields) EncodeURLValOnly(sep string, isURLEncode bool) string {
	var buffer bytes.Buffer
	for i, f := range t {
		if i != 0 {
			buffer.WriteString(sep)
		}
		if isURLEncode {
			buffer.WriteString(url.QueryEscape(fmt.Sprintf("%v", f.Value)))
		} else {
			buffer.WriteString(fmt.Sprintf("%v", f.Value))
		}
	}
	return buffer.String()
}

// EncodeWithFlag encodes TagFields into string, add rules by flag
func (t TagFields) EncodeWithFlag(sep string, flag Flag) string {
	var buffer bytes.Buffer
	isFirstLoop := true
	for _, f := range t {

		val := fmt.Sprintf("%v", f.Value)
		if flag&OmitEmpty != 0 && val == "" {
			continue
		}

		if !isFirstLoop {
			buffer.WriteString(sep)
		}

		// 非ValOnly的写入参数名
		if flag&ValOnly == 0 {
			buffer.WriteString(f.Tag)
			// 非WithoutEquality的写入"="
			if flag&WithoutEquality == 0 {
				buffer.WriteString("=")
			}
		}

		// URLDecode 在TrimSpace之前，因为urlencode过的参数值不含空格
		if flag&URLDecode != 0 {
			val, _ = url.QueryUnescape(val)
		}

		// TrimSpace去除首尾空格
		if flag&TrimSpace != 0 {
			val = strings.TrimSpace(val)
		}

		// URLEncode
		if flag&URLEncode != 0 {
			val = url.QueryEscape(val)
		}

		// DoubleQuotation
		if flag&DoubleQuotation != 0 {
			val = fmt.Sprintf(`"%s"`, val)
		}

		buffer.WriteString(val)
		isFirstLoop = false
	}

	return buffer.String()
}

// GetFields calls GetFieldsByTag with param tagKey="sort"
func GetFields(v interface{}) (fields TagFields) {
	return GetFieldsByTag(v, "sort")
}

// GetFieldsByTag returns tags and its' corresponding values sorted
// by given tagKey in alphabet order.
// nil is return if v is not struct
//
// Example:
//    type User struct{
//		Passwd string `schema:"passwd"`
//		Name   string `schema:"name"`
//		Age    int    `schema:"age"`
//    }
//    u := User{"0000", "alice", 20}
//    fields := GetFieldsSortedByTag(u, "schema")
//    // fields =[{"age", 20}, {"name", "alice"}, {"passwd", "0000"}]
//    GetFieldsSortedByTag("string")
//    // nil is return
//
func GetFieldsByTag(v interface{}, tagKey string) (fields TagFields) {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil
	}

	obj := reflect.ValueOf(v)
	var tfs []TagField
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get(tagKey)
		if t == "" || t == "-" { // escape empty tag and "-" tag
			continue
		}
		name, opts := parseTag(t)
		val := obj.FieldByName(f.Name)
		if opts.Contains("omitempty") && isEmptyValue(val) {
			continue
		}

		tf := TagField{
			Tag:   name,
			Value: val,
		}
		tfs = append(tfs, tf)
	}
	fields = TagFields(tfs)
	sort.Sort(fields)
	return fields
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}

// GetNestFields calls GetNestFieldsByTag with param tagKey="sort"
func GetNestFields(v interface{}) (fields TagFields) {
	return GetNestFieldsByTag(v, "sort")
}

// GetNestFieldsByTag returns tags and its' corresponding values sorted in nested struct
// by given tagKey in alphabet order.
// nil is return if v is not struct
//
/*
Example:
	type User struct {
		Id   string    `schema:"id"`
		Name string `schema:"name"`
		Address
	}

	type Address struct {
		Add string `schema:"add"`
		Res string    `schema:"res"`
		SubAddr
	}

	type SubAddr struct {
		Post string `schema:"post"`
	}

	u := User{
		ID:   "1001",
		Name: "cg",
		Address: Address{
			Add: "local^&",
			Res: "-12",
			SubAddr: SubAddr{
				Post: "comeon. ",
			},
		},
	}

	fields := GetNestFieldsSortedByTag(u, "schema")
	fields = [{"add", "local^&"}, {"id", "1001"}, {"name", "cg"}, {"post", "comeon. "}, {"res", "-12"}]
*/
func GetNestFieldsByTag(v interface{}, tagKey string) (fields TagFields) {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil
	}

	obj := reflect.ValueOf(v)
	tfs := getNestFieldsByTag(rt, obj, tagKey)
	fields = TagFields(tfs)
	sort.Sort(fields)
	return fields
}

func getNestFieldsByTag(rt reflect.Type, obj reflect.Value, tagKey string) []TagField {
	var tfs []TagField
	for i := 0; i < rt.NumField(); i++ {
		if !obj.Field(i).CanInterface() { // if can not be exported
			continue
		}

		if obj.Field(i).Type().Kind() == reflect.Struct { // if nested
			tfs = append(tfs, getNestFieldsByTag(obj.Field(i).Type(), obj.Field(i), tagKey)...)
		} else {
			f := rt.Field(i)
			t := f.Tag.Get(tagKey)
			if t == "" || t == "-" { // escape empty tag and "-" tag
				continue
			}

			name, opts := parseTag(t)
			val := obj.FieldByName(f.Name)
			if opts.Contains("omitempty") && isEmptyValue(val) {
				continue
			}

			tf := TagField{
				Tag:   name,
				Value: val,
			}

			tfs = append(tfs, tf)
		}
	}

	return tfs
}

// FieldValues wraps string slice, and provides Encode method.
// FieldValues can be sort by tag in alphabet order.
type FieldValues []string

// GetValues calls GetValuesByTag with param tagKey="sort"
func GetValues(v interface{}) FieldValues {
	return GetValuesByTag(v, "sort")
}

// GetValuesByTag returns tags and its' corresponding values sorted
// by given tagKey in alphabet order.
// nil is return if v is not struct
//
// Example:
//    type User struct{
//		Passwd string `schema:"passwd"`
//		Name   string `schema:"-"`
//		Age    int    `schema:"omitempty"`
//		Email  string `schema:"email"`
//    }
//    u := User{
//			Passwd: "0000",
//			Name:	"alice",
//			Email:  "maniafish",
//		}
//    values := GetValuesByTag(u, "schema")
//	  values: []string{"0000", "maniafish"}
//    GetFieldsSortedByTag("string")
//    // nil is return
//
func GetValuesByTag(v interface{}, tagKey string) FieldValues {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil
	}

	obj := reflect.ValueOf(v)
	var vals []string
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get(tagKey)
		if t == "" || t == "-" { // escape empty tag and "-" tag
			continue
		}

		_, opts := parseTag(t)
		val := obj.FieldByName(f.Name)
		if opts.Contains("omitempty") && isEmptyValue(val) {
			continue
		}

		vals = append(vals, fmt.Sprintf("%v", val))
	}
	sort.Strings(vals)
	return FieldValues(vals)
}

// Encode encodes TagFields into string by joining [value1, value2, ...]
// seperated by given sep.
//
// Example:
//		given sep=&, t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
//		"val1&val2&...&valn" is return
func (t FieldValues) Encode(sep string) string {
	return t.EncodeURL(sep, false)
}

//EncodeURL is similar to Encode, isURLEncode means whether the value need to be urlencode
func (t FieldValues) EncodeURL(sep string, isURLEncode bool) string {
	var buffer bytes.Buffer
	for i, v := range t {
		if i != 0 {
			buffer.WriteString(sep)
		}

		if isURLEncode {
			buffer.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
		} else {
			buffer.WriteString(fmt.Sprintf("%v", v))
		}
	}

	return buffer.String()
}
