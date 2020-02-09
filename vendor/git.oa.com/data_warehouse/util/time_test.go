package util

import (
	"testing"
)

func TestGetToday(t *testing.T) {
	u := NewDateUtil("2017-03-01")
	d := u.GetToday()
	if d != "2017-03-01" {
		t.Errorf("GetToday, want %s, but get %s", "2017-03-01", d)
	}

	u.SetBaseDate("20170302")
	d = u.GetToday()
	if d != "2017-03-02" {
		t.Errorf("GetToday, want %s, but get %s", "2017-03-02", d)
	}
}

func TestGetYesterday(t *testing.T) {
	u := NewDateUtil("2017-03-01")
	d := u.GetYesterday()
	if d != "2017-02-28" {
		t.Errorf("GetYesterday , want %s, but get %s", "2017-02-28", d)
	}
}

func TestDateAdd(t *testing.T) {
	sd := "2017-02-28"
	ed := DateAdd(sd, 1)
	if ed != "2017-03-01" {
		t.Errorf("DateAdd, want %s, but get %s", "2017-03-01", ed)
	}
}

func TestDateSub(t *testing.T) {
	sd := "2017-02-28"
	ed := DateSub(sd, -1)
	if ed != "2017-03-01" {
		t.Errorf("DateSub, want %s, but get %s", "2017-03-01", ed)
	}
}

func TestTimeRange(t *testing.T) {
	var realResults, expectedResults []string
	realResults = DateRange("20170227", "2017-03-02", FORMAT_FULL_DATE)
	expectedResults = []string{"2017-02-27", "2017-02-28", "2017-03-01", "2017-03-02"}

	if len(realResults) != 4 {
		t.Errorf("DateRange, want %d, but get %d", 4, len(realResults))
	}

	for i := range realResults {
		if realResults[i] != expectedResults[i] {
			t.Errorf("DateRange, want %s, but get %s", expectedResults[i], realResults[i])
		}
	}
}

func TestGetFormatedDateTime(t *testing.T) {
	v := GetFormatedDateTime("20140916", FORMAT_FULL_DATE)
	if v != "2014-09-16" {
		t.Errorf("GetFormatedDateTime failed, want %s, but get %s", "2014-09-16", v)
	}

	v = GetFormatedDateTime("2014-09-16T01:02:03Z", FORMAT_FULL_DATE)
	if v != "2014-09-16" {
		t.Errorf("GetFormatedDateTime failed, want %s, but get %s", "2014-09-16", v)
	}
}
