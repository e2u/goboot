package goboot

import (
	"os"

	logging "github.com/op/go-logging"
)

var (
	Log                       *logging.Logger
	LoggingFormatWithColor    = logging.MustStringFormatter(`%{color}%{time:2006-01-02T15:04:05.9999-07:00} %{id:08x} %{shortfile} %{longfunc} ▶ %{level:-8s} %{color:reset} %{message}`)
	LoggingFormatJSON         = logging.MustStringFormatter(`{"timestamp":"%{time:2006-01-02T15:04:05.9999-07:00}","id":%{id:08x},"filename":"%{shortfile}","func":"%{longfunc}","level":"%{level:s}","msg":"%{message}"}`)
	LoggingFormatWithoutColor = logging.MustStringFormatter(`%{time:2006-01-02T15:04:05.9999-07:00} %{id:08x} %{shortfile} %{longfunc} ▶ %{level:-8s} %{message}`)
	LoggingFormatLogStash     = logging.MustStringFormatter(`{"@timestamp":"%{time:2006-01-02T15:04:05.9999-07:00}","id":%{id:08x},"filename":"%{shortfile}","func":"%{longfunc}","level":"%{level:s}","msg":"%{message}"}`)
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
	format := Config.MustString(IniLogFormat, "plain")
	level := Config.MustString(IniLevel, "DEBUG")
	output := Config.MustString(IniLogOutput, "stdout")

	Log = initLogger(module, format, level, output)
}

func initLogger(module string, format, level, output string) *logging.Logger {
	l := logging.MustGetLogger(module)

	var loggingFormat logging.Formatter

	switch format {
	case "plain":
		loggingFormat = LoggingFormatWithoutColor
	case "plain-color":
		loggingFormat = LoggingFormatWithColor
	case "json":
		loggingFormat = LoggingFormatJSON
	case "logstash":
		loggingFormat = LoggingFormatLogStash
	default:
		loggingFormat = LoggingFormatWithoutColor
	}

	b := getBackend(output)
	formater := logging.NewBackendFormatter(b, loggingFormat)
	backendLeveled := logging.AddModuleLevel(formater)
	lev, err := logging.LogLevel(level)

	if err != nil {
		lev = logging.DEBUG
	}
	backendLeveled.SetLevel(lev, module)
	logging.SetBackend(backendLeveled)
	return l
}

func getBackend(output string) logging.Backend {
	switch output {
	case "off":
		return EmtpyBackend{}
	case "stdout":
		return logging.NewLogBackend(os.Stdout, "", 0)
	case "stderr":
		return logging.NewLogBackend(os.Stdout, "", 0)
	default:
		if out, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0644); err == nil {
			return logging.NewLogBackend(out, "", 0)
		} else {
			return EmtpyBackend{}
		}
	}
}
