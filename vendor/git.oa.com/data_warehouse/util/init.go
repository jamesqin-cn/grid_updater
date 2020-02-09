package util

import (
	"time"
)

func init() {
	go func() {
		for {
			CleanExpiredCache()
			time.Sleep(2 * 1000 * time.Millisecond)
		}
	}()
}
