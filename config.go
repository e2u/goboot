package goboot

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	ini "gopkg.in/ini.v1"
)

type configContext struct {
	*ini.File
	RunModeSection *ini.Section
	DefaultSection *ini.Section
}

func NewConfigWithFile(file, runMode string) *configContext {
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	c := &configContext{
		File: cfg,
		RunModeSection: func() *ini.Section {
			sec, _ := cfg.GetSection(runMode)
			return sec
		}(),
		DefaultSection: func() *ini.Section {
			sec, _ := cfg.GetSection(ini.DEFAULT_SECTION)
			return sec
		}(),
	}
	return c
}

func NewConfigWithoutFile(runMode string) *configContext {
	cfg := ini.Empty()
	envs := []string{ini.DEFAULT_SECTION, "dev", "test", "prod", runMode}
	for _, env := range envs {
		sec, _ := cfg.NewSection(env)
		sec.NewKey("log.output", "stdout")
		sec.NewKey("log.level", "debug")
		sec.NewKey("log.color", "false")
		sec.NewKey("mode.dev", "false")
		sec.NewKey("log.dump.http.request", "true")
		sec.NewKey("log.dump.http.request.body", "true")
		sec.NewKey("log.dump.http.response", "true")
		sec.NewKey("log.dump.http.response.body", "true")
	}
	dsec, _ := cfg.GetSection(ini.DEFAULT_SECTION)
	rsec, _ := cfg.GetSection(runMode)
	return &configContext{
		File:           cfg,
		RunModeSection: rsec,
		DefaultSection: dsec,
	}
}

func (c *configContext) LogLevel() string {
	return c.LogLevel()
}

func (c *configContext) SetModeDev(b bool) {
	c.DefaultSection.Key("mode.dev").SetValue(strconv.FormatBool(b))
}

func (c *configContext) ModeDev() bool {
	return c.MustBool("mode.dev")
}

func (c *configContext) LogDumpHttpRequest() bool {
	return c.MustBool("log.dump.http.request")
}

func (c *configContext) LogDumpHttpRequestBody() bool {
	return c.MustBool("log.dump.http.request.body")
}

func (c *configContext) LogDumpHttpResponse() bool {
	return c.MustBool("log.dump.http.response")
}

func (c *configContext) LogDumpHttpResponseBody() bool {
	return c.MustBool("log.dump.http.response.body")
}

func (c *configContext) mustKeyValue(key string) (*ini.Key, error) {
	switch {
	case c.RunModeSection != nil && c.RunModeSection.HasKey(key):
		return c.RunModeSection.Key(key), nil
	case c.DefaultSection.HasKey(key):
		return c.DefaultSection.Key(key), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid ini key: %s", key))
	}
}

func (c *configContext) MustInt(key string, defaultVal ...int) int {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustInt()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *configContext) MustBool(key string, defaultVal ...bool) bool {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustBool()
	} else if len(defaultVal) == 0 {
		return false
	}
	return defaultVal[0]
}

func (c *configContext) MustDuration(key string, defaultVal ...time.Duration) time.Duration {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustDuration()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *configContext) MustFloat64(key string, defaultVal ...float64) float64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustFloat64()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *configContext) MustString(key string, defaultVal ...string) string {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.String()
	} else if len(defaultVal) == 0 {
		return ""
	}
	return defaultVal[0]
}

func (c *configContext) MustTime(key string, defaultVal ...time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTime()
	} else if len(defaultVal) == 0 {
		t, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
		return t
	}
	return defaultVal[0]
}

func (c *configContext) MustTimeFormat(key, format string, defaultVal ...time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTimeFormat(format)
	} else if len(defaultVal) == 0 {
		t, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
		return t
	}
	return defaultVal[0]
}

func (c *configContext) MustUint(key string, defaultVal ...uint) uint {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *configContext) MustUint64(key string, defaultVal ...uint64) uint64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint64()
	} else if len(defaultVal) == 0 {
		return 0
	}
	return defaultVal[0]
}

func (c *configContext) MustURL(key string, defaultVal ...interface{}) *url.URL {
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

func (c *configContext) MustBase64String(key string, defaultVal ...[]byte) []byte {
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
