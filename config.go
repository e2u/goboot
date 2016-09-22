package goboot

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	ini "gopkg.in/ini.v1"
)

type ConfigContext struct {
	*ini.File
	RunModeSection *ini.Section
	DefaultSection *ini.Section
}

func NewConfig() *ConfigContext {
	return &ConfigContext{}
}

func NewConfigWithFile(file string) *ConfigContext {
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	c := &ConfigContext{
		File: cfg,
		RunModeSection: func() *ini.Section {
			sec, _ := cfg.GetSection(RunMode)
			return sec
		}(),
		DefaultSection: func() *ini.Section {
			sec, err := cfg.GetSection(ini.DEFAULT_SECTION)
			if err != nil {
				panic(err)
			}
			return sec
		}(),
	}
	return c
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

func (c *ConfigContext) MustURL(key string, defaultVal ...*url.URL) *url.URL {
	kv := c.MustString(key)
	if kv == "" && len(defaultVal) == 0 {
		return nil
	} else if kv == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	if u, err := url.Parse(kv); err == nil {
		return u
	} else if len(defaultVal) == 0 {
		return nil
	}
	return defaultVal[0]
}
