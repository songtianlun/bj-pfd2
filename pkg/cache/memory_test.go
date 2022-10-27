package cache

import (
    "testing"
    "time"
)

func TestMemoryGet(t *testing.T) {
    c := initMemoryCache()
    err := memorySet(c, "key", "value", 10*time.Minute)
    if err != nil {
        t.Errorf("memorySet() failed - %v", err)
    }
    v := memoryGet(c, "key")
    if v != "value" {
        t.Error("memoryGet() failed - value not match")
    }
    v = memoryGet(c, "key2")
    if v != "" {
        t.Error("memoryGet() failed - unknown key get value")
    }
}
