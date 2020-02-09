package util

import (
	"time"
)

type CostTime struct {
	startTime time.Time
}

func (t *CostTime) Start() {
	t.startTime = time.Now()
}

func (t *CostTime) GetCostTime() int64 {
	return time.Now().Unix() - t.startTime.Unix()
}
