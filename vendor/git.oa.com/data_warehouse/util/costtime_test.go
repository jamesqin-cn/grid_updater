package util

import (
	"testing"
	"time"
)

func TestGetCostTime(t *testing.T) {
	costtime := CostTime{}
	costtime.Start()
	time.Sleep(1 * 1000 * time.Millisecond)
	sec := costtime.GetCostTime()
	if sec != 1 {
		t.Fatalf("GetCostTime() failed, want:%d, get:%d", 1, sec)
	}

	time.Sleep(1 * 1000 * time.Millisecond)
	sec = costtime.GetCostTime()
	if sec != 2 {
		t.Fatalf("GetCostTime() failed, want:%d, get:%d", 2, sec)
	}
}
