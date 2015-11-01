package service

import (
	"testing"
	"time"
)

func TestLimitToPeriod(t *testing.T) {
	var dur time.Duration
	dur, _ = LimitToPeriod(60)
	t.Log(dur)
	if dur != time.Minute {
		t.Fail()
	}
	dur, _ = LimitToPeriod(120)
	t.Log(dur)
	if dur != (30 * time.Second) {
		t.Fail()
	}
}
