package main

import (
	"database/sql"
	"reflect"
	"strings"
)

func encodeFieldname(k string) string {
	if iskey(k) {
		return k + "_"
	}
	return k
}

func iskey(name string) bool {
	switch name {
	case "abstract", "assert", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "default", "do", "double", "else", "enum", "exports", "extends", "final", "finally", "float", "for", "if", "goto", "implements", "import", "instanceof", "int", "interface", "long", "module", "native", "new", "package", "private", "protected", "public", "requires", "return", "short", "static", "strictfp", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "try", "var", "void", "volatile", "while", "yield":
		return true
	default:
		return false
	}
}

var nullTimeType = reflect.TypeOf(sql.NullTime{})
var nullBoolType = reflect.TypeOf(sql.NullBool{})
var nullByteType = reflect.TypeOf(sql.NullByte{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullInt16Type = reflect.TypeOf(sql.NullInt16{})
var nullInt32Type = reflect.TypeOf(sql.NullInt32{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullStringType = reflect.TypeOf(sql.NullString{})

func goType(rtype reflect.Type, fieldTypeName string) string {

	if strings.Contains(strings.ToUpper(fieldTypeName), "DECIMAL") || strings.Contains(strings.ToUpper(fieldTypeName), "NUMERIC") {
		return "BigDecimal"
	}

	switch {
	case rtype == nullTimeType:
		return "Date"
	case rtype == nullBoolType:
		return "boolean"
	case rtype == nullByteType:
		return "byte"
	case rtype == nullFloat64Type:
		return "double"
	case rtype == nullInt16Type:
		return "short"
	case rtype == nullInt32Type:
		return "int"
	case rtype == nullInt64Type:
		return "long"
	case rtype == nullStringType:
		return "String"
	case rtype.Kind() == reflect.Slice && rtype.Elem().Kind() == reflect.Uint8:
		return "byte[]"
	case rtype.Kind() == reflect.Uint8:
		return "byte"
	case rtype.Kind() == reflect.Uint16:
		return "short"
	case rtype.Kind() == reflect.Uint32:
		return "int"
	case rtype.Kind() == reflect.Uint64:
		return "long"
	case rtype.Kind() == reflect.Uint:
		return "int"
	case rtype.Kind() == reflect.Int8:
		return "byte"
	case rtype.Kind() == reflect.Int16:
		return "short"
	case rtype.Kind() == reflect.Int32:
		return "int"
	case rtype.Kind() == reflect.Int64:
		return "long"
	case rtype.Kind() == reflect.Float32:
		return "float"
	case rtype.Kind() == reflect.Float64:
		return "double"
	case rtype.Kind() == reflect.String:
		return "String"
	case rtype.Kind() == reflect.Bool:
		return "boolean"
	default:
		return "String"
	}
}
