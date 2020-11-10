package datetime

import (
	"time"
)

// MonthRange 获取月的起止日期
func MonthRange(date string) (start, end time.Time, err error) {
	tm1, err := time.Parse(GolangTimeTemplate2, date)
	if err != nil {
		return
	}
	year, month, _ := tm1.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.AddDate(0, 0, 0)
	end = thisMonth.AddDate(0, 1, -1)
	return
}

// MonthDayCount 获取某月的天数
func MonthDayCount(date string) (int, error) {
	start, end, err := MonthRange(date)
	if err != nil {
		return 0, err
	}
	DatetimeSubDay(start, end)
	return DatetimeSubDay(start, end), nil
}

func DatetimeSubDay(start, end time.Time) int {
	return int(end.Sub(start).Hours()/24) + 1
}
