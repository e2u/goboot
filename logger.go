package goboot

import (
	"os"

	logging "github.com/op/go-logging"
)

var (
	Log                       *logging.Logger
	LoggingFormatWithColor    = logging.MustStringFormatter(`%{color}%{time:2006-01-02T15:04:05.9999-07:00} %{id:08x} %{shortfile} %{longfunc} ▶ %{level:-8s} %{color:reset} %{message}`)
	LoggingFormatWithoutColor = logging.MustStringFormatter(`%{time:2006-01-02T15:04:05.9999-07:00} %{id:08x} %{shortfile} %{longfunc} ▶ %{level:-8s} %{message}`)
)

type EmtpyBackend struct{}

func (eb EmtpyBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	return nil
}

func InitLogger() {
	module := Config.MustString("app.name", "unknown")
	InitLoggerWithModule(module)
}

func InitLoggerWithModule(module string) {
	Log = logging.MustGetLogger(module)
	var b logging.Backend
	output := Config.MustString(IniLogOutput, "off")

	switch output {
	case "off":
		b = EmtpyBackend{}
	case "stdout":
		b = logging.NewLogBackend(os.Stdout, "", 0)
	case "stderr":
		b = logging.NewLogBackend(os.Stdout, "", 0)
	default:
		if out, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0644); err == nil {
			b = logging.NewLogBackend(out, "", 0)
		} else {
			b = EmtpyBackend{}
		}
	}

	LoggingFormat := LoggingFormatWithoutColor
	if Config.MustBool(IniLogColoe) {
		LoggingFormat = LoggingFormatWithColor
	}

	formater := logging.NewBackendFormatter(b, LoggingFormat)
	backendLeveled := logging.AddModuleLevel(formater)
	level, err := logging.LogLevel(Config.MustString(IniLevel, "DEBUG"))

	if err != nil {
		level = logging.DEBUG
	}
	backendLeveled.SetLevel(level, module)
	logging.SetBackend(backendLeveled)
}
