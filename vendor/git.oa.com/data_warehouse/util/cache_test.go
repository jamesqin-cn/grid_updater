package util

import (
	"testing"
	"time"
)

func TestCleanCache(t *testing.T) {
	var result interface{}

	CleanCache()
	result = GetFromCache("t01")
	if result != nil {
		t.Errorf("GetFromCache must be nil when cleaned, but get %v", result)
	}
}

func TestGetAndSaveCache(t *testing.T) {
	var result interface{}

	CleanCache()

	for i := 1; i <= 100; i++ {
		SaveToCache("t02", "hello", time.Now().Unix()+3600)
		result = GetFromCache("t02")
		if result != "hello" {
			t.Errorf("GetFromCache(t02) at round %d 1st failed, want %s, but get %v", i, "hello", result)
		}

		SaveToCache("t03", "bye", time.Now().Unix()+3600)
		result = GetFromCache("t03")
		if result != "bye" {
			t.Errorf("GetFromCache(t03) at round %d 1st failed, want %s, but get %v", i, "bye", result)
		}

		result = GetFromCache("t02")
		if result != "hello" {
			t.Errorf("GetFromCache(t02) at round %d 2nd failed, want %s, but get %v", i, "hello", result)
		}

		result = GetFromCache("t03")
		if result != "bye" {
			t.Errorf("GetFromCache(t03) at round %d 2nd failed, want %s, but get %v", i, "bye", result)
		}
	}
}

func TestExpired(t *testing.T) {
	var result interface{}

	CleanCache()

	SaveToCache("t05", "hello", time.Now().Unix()+1)
	result = GetFromCache("t05")
	if result != "hello" {
		t.Errorf("GetFromCache failed, want %s, but get %v", "hello", result)
	}

	time.Sleep(3 * 1000 * time.Millisecond)

	result = GetFromCache("t05")
	if result != nil {
		t.Errorf("GetFromCache failed, want %s, but get %v", "<nil>", result)
	}

}
