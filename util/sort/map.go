package sort

// package sort provides primitives for sorting method in matrix system

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
)

// EncodeMap format map to a string by keys order
// Example:
//		given originMap=
//           {
//             "a": 23,
//             "b": 23.2100,
//             "c": 1234567890,
//             "0": nil,
//             "1": "string *^",
//             "2": true
//           }
//           sep=&, isURLEncode=false
// return "0=&1=string *^&2=true&a=23&b=23.21&c=1234567890"
func EncodeMap(originMap map[string]interface{}, sep string, isURLEncode bool) string {
	keys := getSortedKeys(originMap)
	var buffer bytes.Buffer
	for i, k := range keys {
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(k)
		buffer.WriteString("=")

		v := formatStr(originMap[k])
		if isURLEncode {
			buffer.WriteString(url.QueryEscape(v))
		} else {
			buffer.WriteString(v)
		}
	}
	return buffer.String()
}

// EncodeMapValOnly format map to a string by keys order, without keys
// Example:
//		given originMap=
//           {
//             "a": 23,
//             "b": 23.2100,
//             "c": 1234567890,
//             "0": nil,
//             "1": "string *^",
//             "2": true
//           }
//           sep=&, isURLEncode=false
// return "&string *^&true&23&23.21&1234567890"
func EncodeMapValOnly(originMap map[string]interface{}, sep string, isURLEncode bool) string {
	keys := getSortedKeys(originMap)
	var buffer bytes.Buffer
	for i, k := range keys {
		if i != 0 {
			buffer.WriteString(sep)
		}

		v := formatStr(originMap[k])
		if isURLEncode {
			buffer.WriteString(url.QueryEscape(v))
		} else {
			buffer.WriteString(v)
		}
	}
	return buffer.String()
}

func getSortedKeys(originMap map[string]interface{}) []string {
	i, keys := 0, make([]string, len(originMap))
	for k := range originMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func formatStr(value interface{}) string {
	var ret string
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		ret = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Invalid:
		ret = ""
	default:
		ret = fmt.Sprintf("%v", value)
	}
	return ret
}

// EncodeMapString format map[string]string to string
// Example:
//		given originMap=
//           {
//             "a": 23,
//             "b": 23.2100,
//             "c": 1234567890,
//             "0": nil,
//             "1": "string *^",
//             "2": true
//           }
//           sep=&, isURLEncode=false
// return "0=&1=string *^&2=true&a=23&b=23.2100&c=1234567890"
func EncodeMapString(originMap map[string]string, sep string, isURLEncode bool) string {
	keys := make([]string, 0, len(originMap))
	for k := range originMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for i, k := range keys {
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		v := originMap[k]
		if isURLEncode {
			v = url.QueryEscape(v)
		}
		buffer.WriteString(v)
	}
	return buffer.String()
}

// EncodeMapStringExcept format map[string]string to string, except the keys specified in "except []string"
func EncodeMapStringExcept(originMap map[string]string, except []string, sep string, isURLEncode bool, omitempty bool) string {
	exkeys := make(map[string]struct{})
	for _, k := range except {
		exkeys[k] = struct{}{}
	}
	keys := make([]string, 0, len(originMap))
	for k := range originMap {
		if _, ok := exkeys[k]; !ok {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for i, k := range keys {
		if originMap[k] == "" && omitempty {
			continue
		}
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		v := originMap[k]
		if isURLEncode {
			v = url.QueryEscape(v)
		}
		buffer.WriteString(v)
	}
	return buffer.String()
}
