package util

import (
	"errors"
	"strings"
	"time"
)

const (
	FORMAT_FULL_DATE            = "2006-01-02"           //长日期格式
	FORMAT_FULL_DATE_WITH_SLASH = "2006/01/02"           //长日期格式
	FORMAT_FULL_DATE_COMPACT    = "20060102"             //长日期格式
	FORMAT_DATE                 = "01/02"                //短日期格式
	FORMAT_SHORT_DATE           = "0102"                 //短日期格式
	FORMAT_FULL_TIME            = "15:04:05"             //长时间格式
	FORMAT_SHORT_TIME           = "15:04"                //短时间格式
	FORMAT_FULL_DATETIME        = "2006-01-02 15:04:05"  //日期时间格式
	FORMAT_SHORT_DATETIME       = "01-02 15:04"          //日期时间格式
	FORMAT_DB_DATETIME          = "2006-01-02T15:04:05Z" // MYSQL驱动查询得到的日期格式
)

type DateUtil struct {
	baseDate time.Time
}

func NewDateUtil(d string) *DateUtil {
	u := &DateUtil{}
	u.SetBaseDate(d)
	return u
}

func (u *DateUtil) SetBaseDate(d string) (err error) {
	switch len(d) {
	case 4: //MMDD
		u.baseDate, err = time.Parse(FORMAT_SHORT_DATE, d)
		return
	case 5: // MM/DD
		u.baseDate, err = time.Parse(FORMAT_DATE, d)
		return
	case 8: // YYYYMMDD
		u.baseDate, err = time.Parse(FORMAT_FULL_DATE_COMPACT, d)
		return
	case 10: // YYYY-MM-DD or YYYY/MM/DD
		if strings.Contains(d, "-") {
			u.baseDate, err = time.Parse(FORMAT_FULL_DATE, d)
			return
		}

		if strings.Contains(d, "/") {
			u.baseDate, err = time.Parse(FORMAT_FULL_DATE_WITH_SLASH, d)
			return
		}
		return
	case 11: // MM-DD HH:II
		u.baseDate, err = time.Parse(FORMAT_SHORT_DATETIME, d)
		return
	case 19: // YYYY-MM-DD HH:II:SS
		u.baseDate, err = time.Parse(FORMAT_FULL_DATETIME, d)
		return
	case 20: // YYYY-MM-DDTHH:II:SSZ
		u.baseDate, err = time.Parse(FORMAT_DB_DATETIME, d)
		return
	}

	err = errors.New("Not support this format " + d)
	return
}

func (u *DateUtil) GetToday() (d string) {
	return time.Unix(u.baseDate.Unix(), 0).Format(FORMAT_FULL_DATE)
}

func (u *DateUtil) GetNow() (d string) {
	return time.Unix(u.baseDate.Unix(), 0).Format(FORMAT_FULL_DATETIME)
}

func (u *DateUtil) GetUnix() int64 {
	return u.baseDate.Unix()
}

func (u *DateUtil) GetTimeStruct() time.Time {
	return u.baseDate
}

func (u *DateUtil) GetYesterday() (d string) {
	return u.baseDate.AddDate(0, 0, -1).Format(FORMAT_FULL_DATE)
}

func (u *DateUtil) GetLastWeek() (d string) {
	return u.baseDate.AddDate(0, 0, -7).Format(FORMAT_FULL_DATE)
}

func (u *DateUtil) GetLast2Week() (d string) {
	return u.baseDate.AddDate(0, 0, -14).Format(FORMAT_FULL_DATE)
}

func (u *DateUtil) GetLastMonth() (d string) {
	return u.baseDate.AddDate(0, 0, -31).Format(FORMAT_FULL_DATE)
}

func GetToday() string {
	return time.Now().Format(FORMAT_FULL_DATE)
}

func GetNow() string {
	return time.Now().Format(FORMAT_FULL_DATETIME)
}

func GetFormatedDateTime(d string, format string) string {
	u := NewDateUtil(d)
	return u.GetTimeStruct().Format(format)
}

func GetYesterday() string {
	return time.Now().AddDate(0, 0, -1).Format(FORMAT_FULL_DATE)
}

func DateAdd(d string, intervalDays int) string {
	u := NewDateUtil(d)
	return u.GetTimeStruct().AddDate(0, 0, intervalDays).Format(FORMAT_FULL_DATE)
}

func DateSub(d string, intervalDays int) string {
	return DateAdd(d, -intervalDays)
}

func DateRange(startDate string, endDate string, outputFormat string) (dateRange []string) {
	sdSec := NewDateUtil(startDate).GetUnix()
	edSec := NewDateUtil(endDate).GetUnix()

	if sdSec > edSec {
		sdSec, edSec = edSec, sdSec
	}

	dateRange = []string{}

	i := sdSec
	for i <= edSec {
		dateRange = append(dateRange, time.Unix(i, 0).Format(outputFormat))
		i += 86400
	}

	return
}
