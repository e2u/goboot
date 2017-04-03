package main

import (
	"fmt"

	g "github.com/e2u/goboot"
)

func main() {
	g.Init("dev")
	g.Log.Debug("sample startup...")
	g.OnAppStart(func() error {
		g.Log.Debug("func1....")
		return nil
	})

	g.Log.Debug("MustString: ", g.Config.MustString("sqs.name", "none"))
	g.Log.Debug("MustInt: ", g.Config.MustInt("key.int()", 100))
	g.Log.Debug("MustInt: ", g.Config.MustInt("key.int.noexists", 0))

	g.Log.Debug("debug")
	g.Log.Info("info")
	g.Log.Notice("notice")
	g.Log.Warning("warning")
	g.Log.Error("error")
	g.Log.Critical("critical")
	g.Startup()

	fmt.Println("run mode key values")
	for _, k := range g.Config.RunModeSection.KeyStrings() {
		v, _ := g.Config.RunModeSection.GetKey(k)
		fmt.Println(k, "=>", v)
	}
}
