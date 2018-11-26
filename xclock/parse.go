package xclock

import (
	"github.com/bcampbell/fuzzytime"
	"github.com/pkg/errors"
	"github.com/tkuchiki/parsetime"
	"strings"
	"time"
	"github.com/smcduck/xdsa/xstring"
)

// ParseDatetimeString可以看看https://github.com/olebedev/when，是不是值得参考或者纳入

// FIXME: 无法识别毫秒，比如"2018-11-25 13:21:37.400"
// Parse human-readable date and time string to machine friendly values - unix timestamp
func ParseDatetimeStringFuzz(datetimeString string) (time.Time, error) {
	var tm time.Time
	var tmp string
	var e error
	var dt fuzzytime.DateTime
	ymdChs := []string{
		"年",
		"月",
		"日",
	}
	hmsChs := []string{
		"时",
		"分",
		"秒",
	}
	weeksChs := []string{
		"星期一",
		"星期二",
		"星期三",
		"星期四",
		"星期五",
		"星期六",
		"星期天",
		"星期日",
	}
	monshort := []string{
		"jan",
		"feb",
		"mar",
		"apr",
		"may",
		"june",
		"july",
		"aug",
		"sept",
		"oct",
		"nov",
		"dec",
	}
	monlong := []string{
		"january",
		"february",
		"march",
		"april",
		"may",
		"june",
		"july",
		"agust",
		"september",
		"october",
		"november",
		"december",
	}

	if len(datetimeString) == 0 {
		return tm, errors.New("empty time string")
	}
	if xstring.CountDigit(datetimeString) < 2 {
		return tm, errors.New("invalid time string")
	}

	//fmt.Println("\noriginal string " + datetimeString)

	datetimeString = strings.ToLower(datetimeString)
	datetimeString = strings.Replace(datetimeString, "北京时间", "bjt", -1) // "北京时间" -> "bjt"
	for _, ymd := range ymdChs {                                        // "2016年 12月 20日" or "2016年12月20日" -> "2016-12-20"
		if ymd == "日" {
			datetimeString = strings.Replace(datetimeString, ymd, " ", -1)
		} else {
			datetimeString = strings.Replace(datetimeString, ymd+" ", "-", -1)
			datetimeString = strings.Replace(datetimeString, ymd, "-", -1)
		}
	}
	for _, hms := range hmsChs { // "14时 12分 20秒" or "14时12分20秒" -> "14:12:20"
		if hms == "秒" {
			datetimeString = strings.Replace(datetimeString, hms+" ", " ", -1)
		} else {
			datetimeString = strings.Replace(datetimeString, hms+" ", ":", -1)
			datetimeString = strings.Replace(datetimeString, hms, ":", -1)
		}
	}
	for _, weekday := range weeksChs { // "周一," or "周一" -> ""
		datetimeString = strings.Replace(datetimeString, weekday+",", "", -1)
		datetimeString = strings.Replace(datetimeString, weekday, "", -1)
	}
	datetimeString = strings.Replace(datetimeString, "上午", "", -1) // "上午9:30" -> "9:30 am"
	if strings.Contains(datetimeString, "下午") {
		datetimeString = strings.Replace(datetimeString, "下午", "", -1)
		datetimeString += " pm"
	}
	datetimeString = strings.Replace(datetimeString, "下午", "pm", -1)    // "下午9:30" -> "9:30 pm"
	tmp, e = xstring.ReplaceWithTags(datetimeString, "(", ")", " ", -1) // "12月20日(二)13:06" -> "12月20日 13:06"
	if e == nil {
		datetimeString = tmp
	}
	// 把英文月份缩写前后补上空格, 但是后面带","和"."的不用补空格, 可以识别
	for i, mon := range monshort {
		datetimeString = strings.Replace(datetimeString, mon+".", " "+mon+" ", -1)
		datetimeString = strings.Replace(datetimeString, mon+",", " "+mon+" ", -1)
		// "20dec2016" -> "20 dec 2016"
		if strings.Index(datetimeString, mon) >= 0 && strings.Index(datetimeString, monlong[i]) < 0 {
			datetimeString = strings.Replace(datetimeString, mon, " "+mon+" ", 1)
		}
	}
	datetimeString = strings.Replace(datetimeString, "a.m.", "am", 1)
	datetimeString = strings.Replace(datetimeString, "p.m.", "pm", 1)
	if strings.Count(datetimeString, ".") == 2 {
		datetimeString = strings.Replace(datetimeString, ".", "-", -1)
	}

	//fmt.Println("trimed string " + datetimeString)

	fuzzytime.ExtendYear(2016)
	dt, _, e = fuzzytime.Extract(datetimeString)
	if e == nil {
		//fmt.Println("parsed with fuzzy, " + dt.ISOFormat())
		tm = time.Date(dt.Year(), time.Month(dt.Month()), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, time.FixedZone("UTC", dt.TZOffset()))
		return tm, nil
	} else {

		p, _ := parsetime.NewParseTime()
		tm, e = p.Parse(datetimeString)
		if e == nil {
			//fmt.Println("parsed with parsetime, " + tm.String())
			return tm, nil
		} else {
			return tm, errors.New("can't parse this time string")
		}
	}
}

// strict == true : 以ISO标准严谨识别，仅支持2006-01-02格式
// strict == false : 模糊识别
// 严谨的用于Json的读写，模糊的用于网页解析
func ParseDateString(s string, strict bool) (Date, error) {
	if len(s) == 0 {
		return ZeroDate, nil
	}

	if strict {
		fmts := []string{"2006-01-02"}
		for _, v := range fmts {
			tm, err := time.ParseInLocation(v, s, time.Local)
			if err == nil {
				return TimeToDate(tm), nil
			}
		}
	} else {
		tm, err := ParseDatetimeStringFuzz(s)
		if err == nil {
			return TimeToDate(tm), nil
		}
	}

	return ZeroDate, errors.Errorf("Invalid date string '%s'", s)
}

// strict == true
// 比较严格，用于Json的读写，只支持"YYYY-MM-DD/YYYY-MM-DD"格式
// Date Range的"标准"：http://www.ukoln.ac.uk/metadata/dcmi/date-dccd-odrf/
//
// strict == false
// 比较模糊，用于网页解析，支持如下格式
// 2018-01-02 - 2018-01-03 或者 2018-01-02 ~ 2018-01-03 或者 2017.1.2-2017.1.7
func ParseDateRangeString(s string, strict bool) (DateRange, error) {
	// Check
	ErrDefault := errors.Errorf("Invalid date range string '%s'", s)
	if len(s) == 0 {
		return ZeroDateRange, nil
	}
	if len(s) < 5 {
		return DateRange{}, ErrDefault
	}

	// Parse
	splits := []string{}
	if strict {
		splits = []string{"/"}
	} else {
		splits = []string{" / ", " ~ ", " - ", "/", "~", "-"}
	}
	ss := []string{}
	for _, v := range splits {
		ss = strings.Split(s, v)
		if len(ss) == 2 {
			break
		}
	}
	if len(ss) != 2 {
		return DateRange{}, ErrDefault
	}
	begin, err := ParseDateString(ss[0], strict)
	if err != nil {
		return DateRange{}, err
	}
	if len(ss[0]) > 0 && begin.IsZero() {
		return DateRange{}, ErrDefault
	}
	end, err := ParseDateString(ss[1], strict)
	if err != nil {
		return DateRange{}, err
	}
	if len(ss[1]) > 0 && end.IsZero() {
		return DateRange{}, ErrDefault
	}
	return DateRange{Begin:begin, End:end}, nil
}

func ParseDurationString(durationString string) (time.Duration, error) {
	return time.Second, nil
}
