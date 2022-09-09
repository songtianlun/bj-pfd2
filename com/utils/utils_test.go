package utils

import "testing"

func TestStrToInt64(t *testing.T) {
	var testStr string
	var testInt64 uint64
	testStr = "123"
	testInt64 = 123

	if testInt64 != StrToUInt64(testStr) {
		t.Errorf("StrToUInt64(%s) = %d; want 123", testStr, StrToUInt64(testStr))
	}
}

func TestFloat64ToIntStrRMB(t *testing.T) {
	var testFloat64 float64
	var testStr string
	testFloat64 = 123.456
	testStr = "123"

	resStr := Float64ToIntStrRMB(testFloat64)

	if testStr != resStr {
		t.Errorf("Float64ToStr(%f) = %s; want %s", testFloat64, resStr, testStr)
	}
}

func TestGetTypeString(t *testing.T) {
	if GetTypeString(1) != "int" {
		t.Errorf("GetTypeString(1) = %s; want int", GetTypeString(1))
	}
	if GetTypeString("1") != "string" {
		t.Errorf("GetTypeString(\"1\") = %s; want string", GetTypeString("1"))
	}
	if GetTypeString(1.1) != "float64" {
		t.Errorf("GetTypeString(1.1) = %s; want float64", GetTypeString(1.1))
	}
	if GetTypeString(true) != "bool" {
		t.Errorf("GetTypeString(true) = %s; want bool", GetTypeString(true))
	}
}
