package goboot

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	ini "gopkg.in/ini.v1"
)

type ConfigContext struct {
	*ini.File
	RunModeSection *ini.Section
	DefaultSection *ini.Section
}

func NewConfigWithFile(file, runMode string) *ConfigContext {
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	return newConfigContextWithMode(cfg, runMode)
}

func NewConfigWithoutFile(runMode string) *ConfigContext {
	cfg := ini.Empty()
	envs := []string{ini.DEFAULT_SECTION, "dev", "test", "prod", runMode}
	for _, env := range envs {
		sec, _ := cfg.NewSection(env)
		sec.NewKey(IniLogOutput, "stdout")
		sec.NewKey(IniLevel, "debug")
		sec.NewKey(IniLogFormat, "plain")
		sec.NewKey(IniModeDev, "false")
		sec.NewKey(IniDumpHttpRequest, "true")
		sec.NewKey(IniDumpHttpRequestBody, "true")
		sec.NewKey(IniDumpHttpResponse, "true")
		sec.NewKey(IniDumpHttpResponseBody, "true")
	}

	return newConfigContextWithMode(cfg, runMode)
}

func newConfigContextWithMode(cfg *ini.File, runMode string) *ConfigContext {

	readIncludeFile := func(name string) (*ini.Section, error) {
		nameSplit := strings.Split(name, ":")
		if len(nameSplit) <= 0 {
			return nil, errors.New("Illegal parameters: " + name)
		}
		namePrefix := nameSplit[0]
		path := nameSplit[1]
		if strings.HasPrefix(path, "//") {
			path = path[2:]
		}

		switch namePrefix {
		case "file":
			secCfg, err := ini.Load(path)
			if err != nil {
				return nil, err
			}
			return secCfg.GetSection(ini.DEFAULT_SECTION)
		case "s3":
		default:

		}
		return nil, nil
	}

	processInclude := func(sec *ini.Section) error {
		for _, k := range sec.Keys() {
			kn := strings.TrimSpace(k.Name())
			if !strings.HasPrefix(kn, "@include") {
				continue
			}
			kv := k.MustString("")
			isec, err := readIncludeFile(kv)
			if err != nil {
				panic(err.Error())
			}
			for _, k := range isec.KeyStrings() {
				sec.NewKey(k, isec.Key(k).Value())
			}
		}
		return nil
	}

	return &ConfigContext{
		File: cfg,
		RunModeSection: func() *ini.Section {
			sec, _ := cfg.GetSection(runMode)
			processInclude(sec)
			return sec
		}(),
		DefaultSection: func() *ini.Section {
			sec, _ := cfg.GetSection(ini.DEFAULT_SECTION)
			processInclude(sec)
			return sec
		}(),
	}
}

func (c *ConfigContext) LogLevel() string {
	return c.LogLevel()
}

func (c *ConfigContext) SetModeDev(b bool) {
	c.RunModeSection.Key(IniModeDev).SetValue(strconv.FormatBool(b))
}

func (c *ConfigContext) ModeDev() bool {
	return c.MustBool(IniModeDev)
}

func (c *ConfigContext) LogDumpHttpRequest() bool {
	return c.MustBool(IniDumpHttpRequest)
}

func (c *ConfigContext) LogDumpHttpRequestBody() bool {
	return c.MustBool(IniDumpHttpRequestBody)
}

func (c *ConfigContext) LogDumpHttpResponse() bool {
	return c.MustBool(IniDumpHttpResponse)
}

func (c *ConfigContext) LogDumpHttpResponseBody() bool {
	return c.MustBool(IniDumpHttpResponseBody)
}

func (c *ConfigContext) SetLogDumpHttpRequest(b bool) {
	c.RunModeSection.Key(IniDumpHttpRequest).SetValue(strconv.FormatBool(b))
}

func (c *ConfigContext) SetLogDumpHttpRequestBody(b bool) {
	c.RunModeSection.Key(IniDumpHttpRequestBody).SetValue(strconv.FormatBool(b))
}

func (c *ConfigContext) SetLogDumpHttpResponse(b bool) {
	c.RunModeSection.Key(IniDumpHttpResponse).SetValue(strconv.FormatBool(b))
}

func (c *ConfigContext) SetLogDumpHttpResponseBody(b bool) {
	c.RunModeSection.Key(IniDumpHttpResponseBody).SetValue(strconv.FormatBool(b))
}

func (c *ConfigContext) mustKeyValue(key string) (*ini.Key, error) {
	switch {
	case c.RunModeSection != nil && c.RunModeSection.HasKey(key):
		return c.RunModeSection.Key(key), nil
	case c.DefaultSection.HasKey(key):
		return c.DefaultSection.Key(key), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid ini key: %s", key))
	}
}

func (c *ConfigContext) MustInt(key string, defaultVal ...int) int {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustInt()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustBool(key string, defaultVal ...bool) bool {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustBool()
	} else if len(defaultVal) == 0 {
		return false
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustDuration(key string, defaultVal ...time.Duration) time.Duration {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustDuration()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustFloat64(key string, defaultVal ...float64) float64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustFloat64()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustString(key string, defaultVal ...string) string {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.String()
	} else if len(defaultVal) == 0 {
		return ""
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustTime(key string, defaultVal ...time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTime()
	} else if len(defaultVal) == 0 {
		t, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
		return t
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustTimeFormat(key, format string, defaultVal ...time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTimeFormat(format)
	} else if len(defaultVal) == 0 {
		t, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
		return t
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustUint(key string, defaultVal ...uint) uint {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustUint64(key string, defaultVal ...uint64) uint64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint64()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustURL(key string, defaultVal ...interface{}) *url.URL {
	parseURL := func(v interface{}) *url.URL {
		switch v.(type) {
		case string:
			if u, err := url.Parse(v.(string)); err == nil {
				return u
			}
			return nil
		case *url.URL:
			return v.(*url.URL)
		default:
			return nil
		}
	}

	kv := c.MustString(key)
	if kv == "" && len(defaultVal) == 0 {
		return nil
	} else if kv == "" && len(defaultVal) > 0 {
		return parseURL(defaultVal[0])
	}

	if u, err := url.Parse(kv); err == nil {
		return u
	} else if len(defaultVal) == 0 {
		return nil
	}
	return parseURL(defaultVal[0])
}

func (c *ConfigContext) MustBase64String(key string, defaultVal ...[]byte) []byte {
	kv := c.MustString(key)
	if kv == "" && len(defaultVal) == 0 {
		return nil
	} else if kv == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	if b, err := base64.StdEncoding.DecodeString(kv); err == nil {
		return b
	} else if len(defaultVal) == 0 {
		return nil
	}
	return defaultVal[0]
}

func (c *ConfigContext) MustHexString(key string, defaultVal ...[]byte) []byte {
	kv := c.MustString(key)
	if kv == "" && len(defaultVal) == 0 {
		return nil
	} else if kv == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	if b, err := hex.DecodeString(kv); err == nil {
		return b
	} else if len(defaultVal) == 0 {
		return nil
	}

	return defaultVal[0]
}

func (c *ConfigContext) MustStringArray(key string, sep string, defaultVal ...string) []string {
	kv := c.MustString(key)
	if kv == "" && len(defaultVal) == 0 {
		return nil
	}
	var sr []string
	if as := strings.Split(kv, sep); len(as) > 0 {
		for _, a := range as {
			sr = append(sr, strings.TrimSpace(a))
		}
		return sr
	}
	return strings.Split(defaultVal[0], sep)
}
