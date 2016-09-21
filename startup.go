package goboot

import "sort"

type StartupHook struct {
	order int
	f     func() error
}

type StartupHooks []StartupHook

var startupHooks StartupHooks

func runStartupHooks() {
	var err error
	runFunc := func(f func() error) {
		if err != nil {
			panic(err)
		}
		err = f()
	}

	sort.Sort(startupHooks)
	for _, hook := range startupHooks {
		runFunc(hook.f)
	}
}

func (slice StartupHooks) Len() int {
	return len(slice)
}

func (slice StartupHooks) Less(i, j int) bool {
	return slice[i].order < slice[j].order
}

func (slice StartupHooks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func OnAppStart(f func() error, order ...int) {
	o := 1
	if len(order) > 0 {
		o = order[0]
	}
	startupHooks = append(startupHooks, StartupHook{order: o, f: f})
}
