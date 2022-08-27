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
