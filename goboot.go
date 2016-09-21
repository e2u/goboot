package goboot

import (
	"net/http"
	"strings"
)

var (
	RunMode string
	Config  *ConfigContext
)

func Init(mode string) {
	RunMode = mode
	Config = NewConfigWithFile("conf/app.conf")
	InitLogger()
}

func Startup() {
	runStartupHooks()
	initPprof()
}

func initPprof() {
	go func() {
		ppa := Config.MustString("pprof.addr", "")
		if len(ppa) == 0 {
			return
		}
		pprofUsage :=
			`
Then use the pprof tool to look at the heap profile:

	go tool pprof http://$address$/debug/pprof/heap

Or to look at a 30-second CPU profile:

	go tool pprof http://$address$/debug/pprof/profile
	
Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:

	go tool pprof http://$address$/debug/pprof/block
	
Or to collect a 5-second execution trace:

	wget http://$address$/debug/pprof/trace?seconds=5

To view all available profiles, open http://$address$/debug/pprof/ in your browser.

For a study of the facility in action, visit

https://blog.golang.org/2011/06/profiling-go-programs.html


`

		Log.Info(strings.Replace(pprofUsage, "$address$", ppa, -1))
		Log.Info(http.ListenAndServe(ppa, nil))
	}()
}
