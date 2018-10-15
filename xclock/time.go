package xclock

import "time"

func Sub(t time.Time, d time.Duration) time.Time {
	return t.Add(0 - d)
}

func StringTimeZone(tm time.Time, tz time.Location) string {
	return tm.In(&tz).String()
}

var (
	EpochBeginTime time.Time = EpochSecToTime(0) // 1970-01-01 00:00:00 +0000 UTC
	ZeroTime time.Time = time.Time{} // 0001-01-01 00:00:00 +0000 UTC

	ZeroDate Date = TimeToDate(ZeroTime) // 0001-01-01 00:00:00 +0000 UTC
	ZeroDateRange DateRange = DateRange{Begin:ZeroDate, End:ZeroDate}
)