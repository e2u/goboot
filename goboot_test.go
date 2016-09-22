package goboot

import "testing"

func TestGoBootNew(t *testing.T) {
	OnAppStart(func() error { Log.Debug("001"); return nil }, 1)
	OnAppStart(func() error { Log.Debug("000"); return nil }, 0)
	OnAppStart(func() error { Log.Debug("999"); return nil })
}
