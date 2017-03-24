package goboot

import (
	"bytes"
	"net/url"
	"testing"
	"time"
)

func TestNewConfigWithFileLoadFilePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewConfigWithFile("noexits.conf", "dev")
}

func TestNewConfigWithoutFile(t *testing.T) {
	cfg := NewConfigWithoutFile("dev")

	if cfg.MustString(IniLogOutput) != "stdout" {
		t.Error(IniLogOutput)
	}
}

func TestNewConfigWithFile(t *testing.T) {
	cfg := NewConfigWithFile("config_test.conf", "dev")

	if cfg.MustString("default.app.name") != "DefaultAppName" {
		t.Error("default.app.name")
	}

	if cfg.MustString("app.name", "unknown") != "AppName@dev" {
		t.Error("dev app.name")
	}

	if cfg.MustString("key.noexists", "unknown") != "unknown" {
		t.Error("key.noexits.default")
	}

	if cfg.MustString("key.noexists") != "" {
		t.Error("key.noexists")
	}
	// int test
	if cfg.MustInt("int.100") != 100 {
		t.Error("int.100")
	}

	if cfg.MustInt("int.0") != 0 {
		t.Error("int.0")
	}

	if cfg.MustInt("int.-100") != -100 {
		t.Error("int.-100")
	}

	if cfg.MustInt("int.noexists") != 0 {
		t.Error("int.noexists")
	}

	if cfg.MustInt("int.noexists", 200) != 200 {
		t.Error("int.noexists")
	}

	// bool test
	if cfg.MustBool("bool.true") != true {
		t.Error("bool.true")
	}

	if cfg.MustBool("bool.false") != false {
		t.Error("bool.false")
	}

	if cfg.MustBool("bool.true.1") != true {
		t.Error("bool.true.1")
	}

	if cfg.MustBool("bool.false.0") != false {
		t.Error("bool.false.0")
	}

	if cfg.MustBool("bool.noexits") != false {
		t.Error("bool.noexits")
	}

	if cfg.MustBool("bool.noexists", true) != true {
		t.Error("bool.noexists")
	}

	// time.Duration test
	if cfg.MustDuration("time.duration.1s") != time.Second {
		t.Error("time.duration.1s")
	}

	if cfg.MustDuration("time.duration.1m") != time.Minute {
		t.Error("time.duration.1m")
	}

	if cfg.MustDuration("time.duration.1h5m3s") != time.Hour+5*time.Minute+3*time.Second {
		t.Error("time.duration.1h5m3s")
	}

	if cfg.MustDuration("time.duration.noexists") != 0 {
		t.Error("time.duration.noexists")
	}

	if cfg.MustDuration("time.duration.noexists", time.Hour) != time.Hour {
		t.Error("time.duration.noexists")
	}

	if cfg.MustFloat64("float64.123.45") != 123.45 {
		t.Error("float64.123.45")
	}

	if cfg.MustFloat64("float64.0.123456789") != 0.123456789 {
		t.Error("float64.0.123456789")
	}

	if cfg.MustFloat64("float64.-0.123456789") != -0.123456789 {
		t.Error("float64.-0.123456789")
	}

	if cfg.MustFloat64("float64.noexists") != 0 {
		t.Error("float64.noexists")
	}

	if cfg.MustFloat64("float64.noexists", 0.999999) != 0.999999 {
		t.Error("float64.noexists")
	}

	if cfg.MustString("string.hello") != "hello" {
		t.Error("string.hello")
	}

	if cfg.MustString("string.妳好嗎") != "妳好嗎" {
		t.Error("string.妳好嗎")
	}

	if cfg.MustString("string.你好嗎") != "你好嗎" {
		t.Error("string.你好嗎")
	}

	if cfg.MustString("string.hello.tail.space") != "hello " {
		t.Error("string.hello.tail.space")
	}

	if cfg.MustString("string.hello.head.space") != " hello" {
		t.Error("string.hello.head.space")
	}

	ats := cfg.MustStringArray("string.array", ",")

	if ats[0] != "你好" || ats[1] != "哈哈" || ats[2] != "一直" {
		t.Error("string.array")
	}

	t0, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")

	t1, _ := time.Parse(time.RFC3339, "2016-09-22T10:45:46+08:00")
	if !cfg.MustTime("time.1").Equal(t1) {
		t.Error("time.1")
	}

	t2, _ := time.Parse(time.RFC3339, "2016-09-22T23:45:46-08:00")
	if !cfg.MustTime("time.2").Equal(t2) {
		t.Error("time.2")
	}

	if !cfg.MustTime("time.noexists").Equal(t0) {
		t.Error("time.noexists")
	}

	if !cfg.MustTime("time.noexists", t1).Equal(t1) {
		t.Error("time.noexists")
	}

	t3, _ := time.Parse(time.RFC1123, "Thu, 22 Sep 2016 12:28:02 CST")
	if !cfg.MustTimeFormat("time.format.1", time.RFC1123).Equal(t3) {
		t.Error("time.format.1")
	}

	if !cfg.MustTimeFormat("time.format.noexists", time.RFC1123).Equal(t0) {
		t.Error("time.format.noexists default")
	}

	if !cfg.MustTimeFormat("time.format.noexists", time.RFC1123, t3).Equal(t3) {
		t.Error("time.format.noexists default")
	}

	if cfg.MustUint("uint.100") != 100 {
		t.Error("uint.100")
	}

	if cfg.MustUint("uint.-100") != 0 {
		t.Error("uint.-100")
	}

	if cfg.MustUint("uint.-100", 99) != 99 {
		t.Error("uint.-100")
	}

	if cfg.MustUint64("uint64.100") != 100 {
		t.Error("uint64.100")
	}

	if cfg.MustUint64("uint64.-100") != 0 {
		t.Error("uint64.-100")
	}

	if cfg.MustUint64("uint64.-100", 99) != 99 {
		t.Error("uint64.-100 default")
	}

	url1, _ := url.Parse("https://www.domain.com/path/file.ext")
	if cfg.MustURL("url.1").String() != url1.String() {
		t.Error("url.1")
	}

	url2, _ := url.Parse("http://www.domain.com:7799/path/file.ext")
	if cfg.MustURL("url.2").String() != url2.String() {
		t.Error("url.2")
	}

	url3, _ := url.Parse("h://www.domain.com/path/file.ext")
	if cfg.MustURL("url.3").String() != url3.String() {
		t.Error("url.3")
	}

	if cfg.MustURL("url.noexits") != nil {
		t.Error("url.noexists")
	}

	if cfg.MustURL("url.noexits.default", url1).String() != url1.String() {
		t.Error("url.noexists.default")
	}

	if cfg.MustURL("url.4") != nil {
		t.Error("url.4")
	}

	if cfg.MustURL("url.4", url1).String() != url1.String() {
		t.Error("url.4 default")
	}

	if cfg.MustURL("url.noexits.string", "http://www.domain.com/path.file.ext").String() != "http://www.domain.com/path.file.ext" {
		t.Error("url.noexits.string")
	}

	if cfg.MustURL("url.noexits.string.error", "http://192.168.0.%31:8080/") != nil {
		t.Error("url.noexits.string.error")
	}

	if cfg.MustURL("url.noexits.string.unsuport.type", 10000) != nil {
		t.Error("url.noexits.string.unsuport.type")
	}

	b1 := []byte("hello")
	if bytes.Equal(cfg.MustBase64String("base64.1"), b1) {
		t.Error("base64.1")
	}

	if cfg.MustBase64String("base64.2") != nil {
		t.Error("base64.2")
	}

	if cfg.MustBase64String("base64.noexists") != nil {
		t.Error("base64.noexists")
	}

	if !bytes.Equal(cfg.MustBase64String("base64.2", b1), b1) {
		t.Error("base64.2")
	}

	if !bytes.Equal(cfg.MustBase64String("base64.noexists", b1), b1) {
		t.Error("base64.1.noexists.default")
	}

}
