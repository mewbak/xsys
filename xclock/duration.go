package xclock

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/smcduck/xdsa/xstring"
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

type HumanReadableDuration time.Duration

func ToHumanReadableDuration(d time.Duration) HumanReadableDuration {
	return HumanReadableDuration(d)
}

func (d HumanReadableDuration) ToDuration() time.Duration {
	return time.Duration(d)
}

func (d HumanReadableDuration) String() string {
	return StringHumanReadable(d.ToDuration())
}

func (d HumanReadableDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *HumanReadableDuration) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) <= 1 {
		return errors.Errorf("Invalid json HumanReadableDuration '%s'", s)
	}
	if s[0] != '"' || s[len(s) - 1] != '"' {
		return errors.Errorf("Invalid json HumanReadableDuration '%s'", s)
	}
	s = xstring.RemoveHead(s, 1)
	s = xstring.RemoveTail(s, 1)
	dura, err := time.ParseDuration(s)
	if err != nil {
		*d = HumanReadableDuration(time.Duration(0))
		return err
	}
	*d = HumanReadableDuration(dura)
	return nil
}