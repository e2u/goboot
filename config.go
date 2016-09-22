package goboot

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	ini "gopkg.in/ini.v1"
)

type configContext struct {
	*ini.File
	runModeSection *ini.Section
	defaultSection *ini.Section
}

func NewConfigWithFile(file, runMode string) *configContext {
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	c := &configContext{
		File: cfg,
		runModeSection: func() *ini.Section {
			sec, _ := cfg.GetSection(runMode)
			return sec
		}(),
		defaultSection: func() *ini.Section {
			sec, _ := cfg.GetSection(ini.DEFAULT_SECTION)
			return sec
		}(),
	}
	return c
}

func (c *configContext) mustKeyValue(key string) (*ini.Key, error) {
	switch {
	case c.runModeSection != nil && c.runModeSection.HasKey(key):
		return c.runModeSection.Key(key), nil
	case c.defaultSection.HasKey(key):
		return c.defaultSection.Key(key), nil
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

func (c *configContext) MustURL(key string, defaultVal ...*url.URL) *url.URL {
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
