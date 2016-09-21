package goboot

import "testing"

func TestGoBootNew(t *testing.T) {
	OnAppStart(func() { Logger.Debug("001") }, 1)
	OnAppStart(func() { Logger.Debug("000") }, 0)
	OnAppStart(func() { Logger.Debug("999") })
	gb := New()
	gb.Run()
}
