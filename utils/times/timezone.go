package times

import (
	"time"
)

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone(timezone string) time.Time {
	chinaTimezone, _ := time.LoadLocation(timezone)
	return time.Now().In(chinaTimezone)
}
