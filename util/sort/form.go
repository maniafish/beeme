package sort

// package sort provides primitives for sorting method in gas system

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// EncodeForm encodes url.Values in sorted order
//
// keys in escape map is omitted
// empty string value is omitted when omitempty is true
//
//	Example:
//		values := url.Values{}
//      values.Set("k2","v2")
//      values.Set("k1","v1")
//      values.Set("k3","")
//
//		EncodeForm(values, "&", false, nil) // returns k1=v1&k2=v2&k3=
//		EncodeForm(values, "&", true, nil) // returns k1=v1&k2=v2
//		EncodeForm(values, "&", false, map[string]bool{"k1",true}) // returns k2=v2&k3=
//
func EncodeForm(values url.Values, sep string, omitempty bool, escape map[string]bool) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	isFirstLoop := true
	for _, k := range keys {
		if escape[k] {
			continue
		}
		v := values.Get(k)
		if omitempty && v == "" {
			continue
		}

		if !isFirstLoop {
			buffer.WriteString(sep)
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)
		isFirstLoop = false
	}
	return buffer.String()
}

// EncodeFormWithFlag encodes url.Values in sorted order according to flag
// the supported flags are same to EncodeWithFlag in struct.go
func EncodeFormWithFlag(values url.Values, sep string, flag Flag, escape map[string]bool) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	isFirstLoop := true
	for _, k := range keys {
		if escape[k] {
			continue
		}
		v := values.Get(k)

		// OmitEmpty
		if flag&OmitEmpty != 0 && v == "" {
			continue
		}

		if !isFirstLoop {
			buffer.WriteString(sep)
		}

		// ValOnly
		if flag&ValOnly == 0 {
			buffer.WriteString(k)
			if flag&WithoutEquality == 0 {
				buffer.WriteString("=")
			}
		}

		// URLDecode
		if flag&URLDecode != 0 {
			v, _ = url.QueryUnescape(v)
		}

		// TrimeSpace
		if flag&TrimSpace != 0 {
			v = strings.TrimSpace(v)
		}

		// URLEncode
		if flag&URLEncode != 0 {
			v = url.QueryEscape(v)
		}

		// DoubleQuotation
		if flag&DoubleQuotation != 0 {
			v = fmt.Sprintf(`"%s"`, v)
		}

		buffer.WriteString(v)
		isFirstLoop = false
	}
	return buffer.String()
}

// EncodeFormValOnly encodes url.Values in sorted order
//
// keys in escape map is omitted
// empty string value is omitted while omitempty is true
//
// Example:
//	values := url.Values{}
//	values.Set("k2","v2")
//	values.Set("k1","v1")
//	values.Set("k3","")
//	values.Set("k4","v4")
//
//	EncodeFormValOnly(values, "&", true, nil) // returns v1&v2&v4
//	EncodeFormValOnly(values, "&", false, nil) // returns v1&v2&&v4
//	EncodeFormValOnly(values, "", true, map[string]bool{"k1",true}) // v2v4
//
func EncodeFormValOnly(values url.Values, sep string, omitempty bool, escape map[string]bool) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	isFirstLoop := true
	for _, k := range keys {
		if escape[k] {
			continue
		}
		v := values.Get(k)
		if omitempty && v == "" {
			continue
		}

		if !isFirstLoop {
			buffer.WriteString(sep)
		}
		buffer.WriteString(v)
		isFirstLoop = false
	}
	return buffer.String()
}
