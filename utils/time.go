package utils

import (
	"strconv"
	"time"
)

const (
	DefaultDateFormat     = "2006-01-02"
	DefaultDatetimeFormat = "2006-01-02 15:04:05"
)

// 日期字符串转时间戳
func DateToTime(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation(DefaultDatetimeFormat[:len(date)], date, loc)
	if err != nil {
		return 0
	}
	return stime.Unix()
}

// 时间戳转日期/时间，可传递格式化字符串
func TimeToDate(timestamp int64, params ...string) string {
	if timestamp <= 0 {
		return ""
	}
	format := DefaultDatetimeFormat
	if len(params) >= 1 {
		format = params[0]
	}
	return time.Unix(timestamp, 0).Format(format)
}

// 时间戳转日期int
func TimeToDateInt(timestamp int64, params ...string) int {
	if timestamp <= 0 {
		return 0
	}
	dt := time.Unix(timestamp, 0).Format("20060102")

	n, err := strconv.Atoi(dt)
	if err != nil {
		return 0
	}
	return n
}

// string to time
func StringToTime(str string) time.Time {
	format := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	stime, _ := time.ParseInLocation(format, str, loc)
	return stime
}
