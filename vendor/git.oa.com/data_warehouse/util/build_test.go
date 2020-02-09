package util

import (
	"strings"
	"testing"
)

func TestBuildTime(t *testing.T) {
	buildTime := BuildTime()
	today := GetToday()

	if strings.HasPrefix(buildTime, today) == false {
		t.Fatalf("BuildTime() failed, want:%s HH:II:SS, get:%s", today, buildTime)
	}
}
