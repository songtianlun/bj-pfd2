package cache

import (
    "testing"
    "time"
)

func TestInit(t *testing.T) {
    Init(true, "memory", "", "", 0)
    if manager == nil {
        t.Error("cache Init() failed - manager is nil")
    }
}

func TestSet(t *testing.T) {
    Init(true, "memory", "", "", 0)
    err := Set("key", "value", 10)
    if err != nil {
        t.Errorf("cache Set() failed - %v", err)
    }
}

func TestGet(t *testing.T) {
    Init(true, "memory", "", "", 0)
    err := Set("key", "value", 10*time.Minute)
    if err != nil {
        t.Errorf("cache Set() failed - %v", err)
    }
    v := Get("key")
    if v != "value" {
        t.Errorf("cache Get() failed \n\texpected: %v\n\tgot: %v", "value", v)
    }
    v = Get("key2")
    if v != "" {
        t.Error("cache Get() failed - unknown key get value")
    }
}
