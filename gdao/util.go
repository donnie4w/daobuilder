package main

import (
	"database/sql"
	"reflect"
	"regexp"
)

func CheckReserveKey(k string) string {
	if b, _ := regexp.MatchString("break|default|func|interface|select|case|defer|go|map|struct|chan|else|goto|package|switch|const|fallthrough|if|range|type|continue|for|import|return|var", k); b {
		return k + "_0"
	}
	return k
}

var nullTimeType = reflect.TypeOf(sql.NullTime{})
var nullBoolType = reflect.TypeOf(sql.NullBool{})
var nullByteType = reflect.TypeOf(sql.NullByte{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullInt16Type = reflect.TypeOf(sql.NullInt16{})
var nullInt32Type = reflect.TypeOf(sql.NullInt32{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullStringType = reflect.TypeOf(sql.NullString{})

func goPtrType(rtype reflect.Type) string {
	switch {
	case rtype == nullTimeType:
		return "time.Time"
	case rtype == nullBoolType:
		return "*bool"
	case rtype == nullByteType:
		return "*byte"
	case rtype == nullFloat64Type:
		return "*float64"
	case rtype == nullInt16Type:
		return "*int16"
	case rtype == nullInt32Type:
		return "*int32"
	case rtype == nullInt64Type:
		return "*int64"
	case rtype == nullStringType:
		return "*string"
	case rtype.Kind() == reflect.Slice && rtype.Elem().Kind() == reflect.Uint8:
		return "[]byte"
	case rtype.Kind() == reflect.Uint8:
		return "*uint8"
	case rtype.Kind() == reflect.Uint16:
		return "*uint16"
	case rtype.Kind() == reflect.Uint32:
		return "*uint32"
	case rtype.Kind() == reflect.Uint64:
		return "*uint64"
	case rtype.Kind() == reflect.Uint:
		return "*uint"
	case rtype.Kind() == reflect.Int8:
		return "*int8"
	case rtype.Kind() == reflect.Int16:
		return "*int16"
	case rtype.Kind() == reflect.Int32:
		return "*int32"
	case rtype.Kind() == reflect.Int64:
		return "*int64"
	case rtype.Kind() == reflect.Float32:
		return "*float32"
	case rtype.Kind() == reflect.Float64:
		return "*float64"
	case rtype.Kind() == reflect.String:
		return "*string"
	case rtype.Kind() == reflect.Bool:
		return "*bool"
	default:
		return "*string"
	}
}

func goType(rtype reflect.Type) string {
	switch {
	case rtype == nullTimeType:
		return "time.Time"
	case rtype == nullBoolType:
		return "bool"
	case rtype == nullByteType:
		return "byte"
	case rtype == nullFloat64Type:
		return "float64"
	case rtype == nullInt16Type:
		return "int16"
	case rtype == nullInt32Type:
		return "int32"
	case rtype == nullInt64Type:
		return "int64"
	case rtype == nullStringType:
		return "string"
	case rtype.Kind() == reflect.Slice && rtype.Elem().Kind() == reflect.Uint8:
		return "[]byte"
	case rtype.Kind() == reflect.Uint8:
		return "uint8"
	case rtype.Kind() == reflect.Uint16:
		return "uint16"
	case rtype.Kind() == reflect.Uint32:
		return "uint32"
	case rtype.Kind() == reflect.Uint64:
		return "uint64"
	case rtype.Kind() == reflect.Uint:
		return "uint"
	case rtype.Kind() == reflect.Int8:
		return "int8"
	case rtype.Kind() == reflect.Int16:
		return "int16"
	case rtype.Kind() == reflect.Int32:
		return "int32"
	case rtype.Kind() == reflect.Int64:
		return "int64"
	case rtype.Kind() == reflect.Float32:
		return "float32"
	case rtype.Kind() == reflect.Float64:
		return "float64"
	case rtype.Kind() == reflect.String:
		return "string"
	case rtype.Kind() == reflect.Bool:
		return "bool"
	default:
		return "string"
	}
}

func mustPtr(rtype reflect.Type) bool {
	switch {
	case rtype == nullTimeType:
		return false
	case rtype.Kind() == reflect.Slice && rtype.Elem().Kind() == reflect.Uint8:
		return false
	default:
		return true
	}
}
