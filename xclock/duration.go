package xclock

import (
	"time"
	"math"
	"github.com/hako/durafmt"
)

const (
	Day        = 24 * time.Hour
	Week       = 7 * Day
	MonthFUZZY = 30 * Day
	Year365    = 365 * Day
)

// 计算两个时间直接的间隔天数，不足一天的算一天称为biggerDays，不足一天不计算称为smallerDays
func DaysBetween(since, to time.Time) (exactDays float64, biggerDays, smallerDays int) {
	days := to.Sub(since).Hours() / float64(24)
	return days, int(math.Ceil(days)), int(math.Floor(days))
}

func StringHumanReadable(duration time.Duration) string {
	return durafmt.Parse(duration).String()
}