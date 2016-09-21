package goboot

import (
	"errors"
	"fmt"
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
	case c.RunModeSection.HasKey(key):
		return c.RunModeSection.Key(key), nil
	case c.DefaultSection.HasKey(key):
		return c.DefaultSection.Key(key), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid ini key: %s", key))
	}
}

func (c *ConfigContext) MustInt(key string, defaultVal int) int {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustInt()
	}
	return defaultVal
}

func (c *ConfigContext) MustBool(key string, defaultVal bool) bool {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustBool()
	}
	return defaultVal
}

func (c *ConfigContext) MustDuration(key string, defaultVal time.Duration) time.Duration {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustDuration()
	}
	return defaultVal
}

func (c *ConfigContext) MustFloat64(key string, defaultVal float64) float64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustFloat64()
	}
	return defaultVal
}

func (c *ConfigContext) MustString(key, defaultVal string) string {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.String()
	}
	return defaultVal
}

func (c *ConfigContext) MustTime(key string, defaultVal time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTime()
	}
	return defaultVal
}

func (c *ConfigContext) MustTimeFormat(key, format string, defaultVal time.Time) time.Time {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustTimeFormat(format)
	}
	return defaultVal
}

func (c *ConfigContext) MustUint(key string, defaultVal uint) uint {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint()
	}
	return defaultVal
}

func (c *ConfigContext) MustUint64(key string, defaultVal uint64) uint64 {
	if v, err := c.mustKeyValue(key); err == nil {
		return v.MustUint64()
	}
	return defaultVal
}
