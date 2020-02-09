package util

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

type RecordSlice struct {
	data         []map[string]string
	key          string
	compareAsStr bool
	isAsc        bool
}

func NewRecordSlice(records []map[string]string, key string, compareAsStr bool, isAsc bool) *RecordSlice {
	return &RecordSlice{records, key, compareAsStr, isAsc}
}

func (r *RecordSlice) Len() int {
	return len(r.data)
}
func (r *RecordSlice) Swap(i, j int) {
	r.data[i], r.data[j] = r.data[j], r.data[i]
}
func (r *RecordSlice) Less(i, j int) bool {
	var isLess bool
	if r.compareAsStr {
		if strings.Compare(r.data[i][r.key], r.data[j][r.key]) == -1 {
			isLess = true
		} else {
			isLess = false
		}
	} else {
		f1 := AnyToFloat(r.data[i][r.key])
		f2 := AnyToFloat(r.data[j][r.key])
		if math.Abs(f1-f2) < 0.000001 {
			isLess = false
		} else if f1 < f2 {
			isLess = true
		} else {
			isLess = false
		}
	}

	if r.isAsc {
		return isLess
	}
	return !isLess
}

func TableLeftJoin(priKey string, tables ...[]map[string]string) []map[string]string {
	newTable := make([]map[string]string, 0, len(tables[0]))
	baseTable := tables[0]
	for _, row := range baseTable {
		if findVal, ok := row[priKey]; ok {
			mergeMapRes := MergeMap(row)
			for _, nextTable := range tables[1:] {
				findMap := TableLookUpRow(nextTable, priKey, findVal)
				if len(findMap) != 0 {
					mergeMapRes = MergeMap(mergeMapRes, findMap)
				}
			}
			if len(mergeMapRes) != 0 {
				newTable = append(newTable, mergeMapRes)
			}
		}
	}
	return TableFillMissingColumn(newTable)
}

func TableFillMissingColumn(table []map[string]string) []map[string]string {
	newTable := make([]map[string]string, 0)
	columns := make(map[string]bool, 0)
	for _, row := range table {
		for k := range row {
			columns[k] = true
		}
	}

	for _, row := range table {
		newRow := make(map[string]string, len(columns))
		for col := range columns {
			if _, ok := row[col]; ok {
				newRow[col] = row[col]
			} else {
				newRow[col] = ""
			}
		}
		newTable = append(newTable, newRow)
	}

	return newTable
}

func MakeDateRangeTable(startDate string, endDate string, priKey string, isAsc bool) (table []map[string]string) {
	reverseRange := func(rs []string) {
		len := len(rs)
		for i := 0; i < len/2; i++ {
			rs[i], rs[len-i-1] = rs[len-i-1], rs[i]
		}
	}

	dateRange := DateRange(startDate, endDate, FORMAT_FULL_DATE)
	if isAsc == false {
		reverseRange(dateRange)
	}

	for _, d := range dateRange {
		table = append(table, map[string]string{priKey: d})
	}
	return
}

func ReverseRecords(records []map[string]string) (newRecords []map[string]string) {
	newRecords = make([]map[string]string, len(records))
	len := len(records)
	for i, v := range records {
		newRecords[len-i-1] = v
	}
	return
}

func SortRecords(records []map[string]string, key string, compareAsStr bool, isAsc bool) {
	sort.Sort(NewRecordSlice(records, key, compareAsStr, isAsc))
}

func MergeMap(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, row := range maps {
		for k, v := range row {
			result[k] = v
		}
	}
	return result
}

func TableLookUpVal(table []map[string]string, findCol string, findVal string, returnCol string, defaultVal string) string {
	for _, row := range table {
		if v, ok := row[findCol]; ok {
			if v == findVal {
				if v, ok := row[returnCol]; ok {
					return v
				}
				return defaultVal
			}
		}
	}
	return defaultVal
}

func TableLookUpRow(table []map[string]string, findCol string, findVal string) map[string]string {
	for _, row := range table {
		if v, ok := row[findCol]; ok {
			if v == findVal {
				return row
			}
		}
	}
	return nil
}

func SumInTable(table []map[string]string, findCol string) (sum float64) {
	sum = 0
	for _, row := range table {
		if v, ok := row[findCol]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				sum += f
			}
		}
	}
	return
}
