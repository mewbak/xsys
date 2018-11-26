package xclock

import (
	"time"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smcduck/xdsa/xstring"
)

func TodayString() string {
	return time.Now().Format("2006-01-02")
}

type Date time.Time

func TimeToDate(tm time.Time) Date {
	return Date(tm)
}

// check if date is valid
// invalid date example: 2018-2-30
func DateValid(year, month, day int) bool {
	if month <= 0 || month >= 13 {
		return false
	}
	if day <= 0 || day >= 32 {
		return false
	}
	tm := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if tm.Year() != year || tm.Month() != time.Month(month) || tm.Day() != day {
		return false
	}
	return true
}

func NewDate(year, month, day int, tz time.Location) (Date, error) {
	if !DateValid(year, month, day) {
		return Date{}, errors.Errorf("Invalid date input %d-%d-%d", year, month, day)
	}
	return Date(time.Date(year, time.Month(month), day, 0, 0, 0, 0, &tz)), nil
}

func (d *Date) Year() int {
	return time.Time(*d).Year()
}

func (d *Date) Month() time.Month {
	return time.Time(*d).Month()
}

func (d *Date) Day() int {
	return time.Time(*d).Day()
}

func (d *Date) Equal(cmp Date) bool {
	srcTime := d.ToTime()
	cmpTime := cmp.ToTime()
	return srcTime.UTC().Unix() == cmpTime.UTC().Unix() &&
		srcTime.UTC().UnixNano() == cmpTime.UTC().UnixNano()
}

func (d *Date) IsZero() bool {
	return d.Equal(ZeroDate)
}

func (d Date) String() string {
	if d.ToTime().IsZero() {
		return ""
	}
	return d.StringYYYY_MM_DD()
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) <= 1 {
		return errors.Errorf("Invalid json date '%s'", s)
	}
	if s[0] != '"' || s[len(s) - 1] != '"' {
		return errors.Errorf("Invalid json date '%s'", s)
	}
	s = xstring.RemoveHead(s, 1)
	s = xstring.RemoveTail(s, 1)
	dt, err := ParseDateString(s, true)
	if err != nil {
		*d = ZeroDate
		return err
	}
	*d = dt
	return nil
}

// yyyymmdd
// Notice:
// 如果你写成了fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)，编译也能通过
// 但是返回结果却是很大很大的数字，因为它们代表函数地址
func (d *Date) StringYYYYMMDD() string {
	return d.ToTime().Format("20060102")
}

// yyyy-mm-dd
func (d *Date) StringYYYY_MM_DD() string {
	return d.ToTime().Format("2006-01-02")
}

// 以time.Time相同的格式输出字符串
func (d *Date) StringTime() string {
	return d.ToTime().String()
}


func (d *Date) ToTime() time.Time {
	return time.Time(*d)
}

type DateRange struct {
	Begin Date
	End   Date
}

func (dr DateRange) String() string {
	if dr.Begin.IsZero() && dr.End.IsZero() {
		return ""
	}
	return dr.Begin.String() + "/" + dr.End.String()
}

func (dr DateRange) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", dr.String())), nil
}

func (dr *DateRange) UnmarshalJSON(b []byte) error {
	// Init
	dr.Begin = ZeroDate
	dr.End = ZeroDate

	// Remove '"'
	s := string(b)
	ErrDefault := errors.Errorf("Invalid json date range '%s'", s)
	if len(s) <= 1 {
		return ErrDefault
	}
	if s[0] != '"' || s[len(s) - 1] != '"' {
		return ErrDefault
	}
	s = xstring.RemoveHead(s, 1)
	s = xstring.RemoveTail(s, 1)

	// Parse
	res, err := ParseDateRangeString(s, true)
	if err != nil {
		return ErrDefault
	}
	*dr = res
	return nil
}

func (dr *DateRange) IsZero() bool {
	return dr.Begin.IsZero() && dr.End.IsZero()
}