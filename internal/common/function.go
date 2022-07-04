package common

import (
	"encoding/json"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Function struct {
}

func (f *Function) SelectExists(db *gorm.DB) bool {
	var count int64
	db.Limit(1).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func (f *Function) NowDate(format ...string) string {
	layout := "Y-m-d H:i:s"
	if len(format) > 0 {
		layout = format[0]
	}

	layout = strings.ReplaceAll(layout, "Y", "2006")
	layout = strings.ReplaceAll(layout, "m", "01")
	layout = strings.ReplaceAll(layout, "d", "02")
	layout = strings.ReplaceAll(layout, "H", "15")
	layout = strings.ReplaceAll(layout, "i", "04")
	layout = strings.ReplaceAll(layout, "s", "05")

	return time.Now().Format(layout)
}

func (f *Function) JSONEncode(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

func (f *Function) JSONDecode(value string) interface{} {
	var v interface{}
	_ = json.Unmarshal([]byte(value), &v)
	return v
}
