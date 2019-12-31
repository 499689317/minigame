package timer

import (
	"github.com/astaxie/beego"
	"time"
)

const (
	OneDay_SecCnt = 24 * 3600
	INT_MAX       = 0x7FFFFFFF
)

func IsToday(day int) bool { return time.Now().Day() == day }
func WeekInYear() int {
	_, ret := time.Now().ISOWeek()
	return ret
}
func GetTodayBeginSec() int64 {
	now := time.Now()
	todayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return todayTime.Unix()
}
func GetTodayEndSec() int64 {
	return GetTodayBeginSec() + OneDay_SecCnt
}
func GetTodayRunSec() int64 {
	now := time.Now()
	return int64(now.Hour()*3600 + now.Minute()*60 + now.Second())
}
func GetTodayLeftSec() int64 {
	return OneDay_SecCnt - GetTodayRunSec()
}

// 时间戳--日期
const g_time_layout = "2006/01/02 15:04:05"
// time.Local表示当前服务器时区
func Str2Time(date string) int64 {
	if v, err := time.ParseInLocation(g_time_layout, date, time.Local); err == nil {
		return v.Unix()
	} else {
		beego.Error(err.Error())
		return 0
	}
}
func Time2Str(sec int64) string { return time.Unix(sec, 0).Format(g_time_layout) }
