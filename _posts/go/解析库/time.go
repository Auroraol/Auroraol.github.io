package adapter

import (
	"strconv"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

var loc *time.Location

func Init() error {
	cnLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	loc = cnLoc
	return nil
}

func LocalTimeNow() time.Time {
	return time.Now().In(loc)
}

// "2024-09-01 00:00:00" 转成日期类型
func StringToTime(stringTime string) (time.Time, error) {
	time, err := time.ParseInLocation(TimeLayout, stringTime, loc)
	if err != nil {
		return LocalTimeNow(), err
	}
	return time, nil
}

// 1725120000 转成日期类型
func TimestampToTime(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return LocalTimeNow(), err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

// 2024-09-01 00:00:00 转成 1725120000
func StringTimeToTimestamp(stringTime string) (int64, error) {
	timeInt64, err := time.ParseInLocation(TimeLayout, stringTime, loc)
	if err != nil {
		return 0, err
	}
	timeUnix := timeInt64.Unix()
	return timeUnix, nil
}

// 1725120000 转成 2024-09-01 00:00:00
func TimestampToStringTime(timestamp int64) string {
	return time.Unix(timestamp, 0).In(loc).Format(TimeLayout)
}
