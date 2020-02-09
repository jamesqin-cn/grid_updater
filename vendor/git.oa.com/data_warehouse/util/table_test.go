package util

import (
	"fmt"
	"math"
	"testing"
)

func TestSumInTable(t *testing.T) {
	arr1 := []map[string]string{
		map[string]string{
			"date":   "2018-01-01",
			"amount": "1.1",
		},
		map[string]string{
			"date":   "2018-01-02",
			"amount": "1.2",
		},
		map[string]string{
			"date":   "2018-01-03",
			"amount": "1.3",
		},
		map[string]string{
			"date": "2018-01-04",
		},
	}
	get1 := SumInTable(arr1, "amount")
	want1 := 3.6
	if math.Abs(want1-get1) >= 0.000001 {
		t.Errorf("SumInArray() failed, want %f, bug get %f", want1, get1)
	}
}

func TestTableLookUp(t *testing.T) {
	arr1 := []map[string]string{
		map[string]string{
			"date":   "2018-01-01",
			"amount": "1.1",
		},
		map[string]string{
			"date":   "2018-01-02",
			"amount": "1.2",
		},
		map[string]string{
			"date":   "2018-01-03",
			"amount": "1.3",
		},
		map[string]string{
			"date": "2018-01-04",
		},
	}
	v := TableLookUpVal(arr1, "date", "2018-01-02", "amount", "0")
	if v != "1.2" {
		t.Errorf("ArrayLookUp() failed, want %s, bug get %s", "1.2", v)
	}

	v = TableLookUpVal(arr1, "date", "2018-01-09", "amount", "0")
	if v != "0" {
		t.Errorf("ArrayLookUp() failed, want %s, bug get %s", "0", v)
	}
}

func TestTableLeftJoin(t *testing.T) {
	arr1 := []map[string]string{
		map[string]string{
			"date": "2018-01-01",
			"col1": "1-1",
			"col2": "1-2",
		},
		map[string]string{
			"date": "2018-01-02",
			"col1": "2-1",
		},
		map[string]string{
			"date": "2018-01-03",
			"col1": "3-1",
			"col2": "3-2",
		},
	}

	arr2 := []map[string]string{
		map[string]string{
			"date": "2018-01-01",
			"col3": "1-3",
			"col4": "1-4",
			"col2": "replace with 1-2",
		},
		map[string]string{
			"date": "2018-01-02",
			"col3": "2-3",
			"col4": "2-4",
		},
		map[string]string{
			"date": "2018-01-03",
			"col3": "3-3",
			"col4": "3-4",
		},
		map[string]string{
			"date": "2018-01-04",
			"col3": "cc3",
			"col4": "dd3",
		},
	}
	arr3 := []map[string]string{
		map[string]string{
			"date": "2018-01-01",
			"col5": "1-5",
		},
	}

	res := TableLeftJoin("date", arr1, arr2, arr3)
	fmt.Println(res)
}

func TestMakeDateRangeTable(t *testing.T) {
	table := MakeDateRangeTable("2018-01-01", "2018-01-05", "date", true)
	fmt.Println(table)

	table = MakeDateRangeTable("2018-01-05", "2018-01-01", "date", true)
	fmt.Println(table)

	table = MakeDateRangeTable("2018-01-01", "2018-01-05", "date", false)
	fmt.Println(table)

	table = MakeDateRangeTable("2018-01-05", "2018-01-01", "date", false)
	fmt.Println(table)
}

func TestReverseRecords(t *testing.T) {
	arr1 := []map[string]string{
		map[string]string{
			"date": "2018-01-01",
			"col1": "1-1",
			"col2": "1-2",
		},
		map[string]string{
			"date": "2018-01-02",
			"col1": "2-1",
		},
		map[string]string{
			"date": "2018-01-03",
			"col1": "3-1",
			"col2": "3-2",
		},
	}

	arr1 = ReverseRecords(arr1)
	firstRecordKey := arr1[0]["date"]
	if firstRecordKey != "2018-01-03" {
		t.Errorf("ReverseRecords 1st failed, want %s, but get %s", "2018-01-03", firstRecordKey)
	}

	arr1 = ReverseRecords(arr1)
	firstRecordKey = arr1[0]["date"]
	if firstRecordKey != "2018-01-01" {
		t.Errorf("ReverseRecords 2nd failed, want %s, but get %s", "2018-01-01", firstRecordKey)
	}
}

func TestSortRecords(t *testing.T) {
	arr1 := []map[string]string{
		map[string]string{
			"id":   "5",
			"date": "2018-01-04",
			"col1": "4-1",
		},
		map[string]string{
			"id":   "3",
			"date": "2018-01-02",
			"col1": "2-1",
		},
		map[string]string{
			"id":   "10",
			"date": "2018-01-03",
			"col1": "3-1",
			"col2": "3-2",
		},
		map[string]string{
			"id":   "8",
			"date": "2018-01-01",
			"col1": "1-1",
			"col2": "1-2",
		},
	}

	SortRecords(arr1, "date", true, true)
	firstRecordKey := arr1[0]["date"]
	if firstRecordKey != "2018-01-01" {
		t.Errorf("SortRecords 1st failed, want %s, but get %s", "2018-01-01", firstRecordKey)
	}

	SortRecords(arr1, "date", true, false)
	firstRecordKey = arr1[0]["date"]
	if firstRecordKey != "2018-01-04" {
		t.Errorf("SortRecords 2nd failed, want %s, but get %s", "2018-01-04", firstRecordKey)
	}

	SortRecords(arr1, "id", false, true)
	firstRecordKey = arr1[0]["date"]
	if firstRecordKey != "2018-01-02" {
		t.Errorf("SortRecords 3rd failed, want %s, but get %s", "2018-01-02", firstRecordKey)
	}

	SortRecords(arr1, "id", false, false)
	firstRecordKey = arr1[0]["date"]
	if firstRecordKey != "2018-01-03" {
		t.Errorf("SortRecords 4times failed, want %s, but get %s", "2018-01-03", firstRecordKey)
	}
}
