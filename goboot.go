package goboot

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

var (
	runMode string
	Config  *ConfigContext
)

func Init(mode ...string) {
	if len(mode) == 0 {
		runMode = "auto"
	} else {
		runMode = mode[0]
	}

	if _, err := os.Stat("conf/app.conf"); os.IsNotExist(err) {
		Config = NewConfigWithoutFile(runMode)
	} else {
		Config = NewConfigWithFile("conf/app.conf", runMode)
	}

	InitLogger()
}

func InitWithModeAndFile(mode, file string) {
	runMode = mode
	Config = NewConfigWithFile(file, runMode)
	InitLogger()
}

func Startup() {
	runStartupHooks()
	initPprof()
}

func RunMode() string {
	return runMode
}
func initPprof() {
	go func() {
		ppa := Config.MustString("pprof.addr", "")
		if len(ppa) == 0 {
			return
		}

		pprofMux := http.DefaultServeMux
		http.DefaultServeMux = http.NewServeMux()

		pprofUsage :=
			`
Then use the pprof tool to look at the heap profile:

	go tool pprof http://$address$/debug/pprof/heap

Or to look at a 30-second CPU profile:

	go tool pprof http://$address$/debug/pprof/profile
	
Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:

	go tool pprof http://$address$/debug/pprof/block
	
Or to collect a 5-second execution trace:

	wget -O trace.out http://$address$/debug/pprof/trace?seconds=5

    go tool trace http://$address$/debug/pprof/trace

To view all available profiles, open http://$address$/debug/pprof/ in your browser.

For a study of the facility in action, visit

https://blog.golang.org/2011/06/profiling-go-programs.html


`

		Log.Info(strings.Replace(pprofUsage, "$address$", ppa, -1))
		Log.Info(http.ListenAndServe(ppa, pprofMux))
	}()
}
