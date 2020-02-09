package util

import (
	"testing"
)

func TestLoadDBConfig(t *testing.T) {
	opts, err := LoadDBConfig("online_ro")
	if err != nil {
		t.Fatalf("LoadDBConfig() failed")
	}
	if opts.MaxIdleConns != 10 {
		t.Fatalf("LoadDBConfig() failed, MaxIdleConns want:%d, get:%d", 10, opts.MaxIdleConns)
	}
}
