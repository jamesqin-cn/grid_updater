package util

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

//BuildLcImgUrl line chart
func BuildLcImgUrl(chs string, chtt string, chmixer string, chm string, chdl []string, chxl []string, chd [][]float64) (chartUrl string, err error) {
	if chm != "N" && chm != "o" && chm != "-" {
		err = errors.New("chm do not suport. usage:chm=[N|o|-]")
		return
	}

	if len(chdl) != len(chd) {
		err = errors.New("length of chdl(chart data legend) must equals to the length of chd(chart data)")
		return
	}

	if len(chxl) != len(chd[0]) {
		err = errors.New("length of chxl(chart axis label) must equals to the length of chd(chart data)")
		return
	}

	v := url.Values{}
	v.Set("cht", "lc")
	v.Set("chs", chs)
	v.Set("chtt", url.QueryEscape(chtt))
	v.Set("chmixer", chmixer)
	v.Set("chm", chm)
	v.Set("chdl", strings.Join(chdl, "|"))
	v.Set("chxl", "0:|"+strings.Join(chxl, "|"))

	var data []string
	for _, row := range chd {
		var newRow []string
		for _, v := range row {
			newRow = append(newRow, strconv.FormatFloat(v, 'f', 2, 64))
		}
		data = append(data, strings.Join(newRow, ","))
	}
	v.Set("chd", "t:"+strings.Join(data, "|"))

	chartUrl = "http://chart.oa.com/chart.php?" + v.Encode()
	return
}

//BuildPcImgURL pie chart
func BuildPcImgURL(chs string, chtt string, chl []string, chd []string) (chartURL string, err error) {
	v := url.Values{}
	v.Set("cht", "p")
	v.Set("chs", chs)
	v.Set("chtt", url.QueryEscape(chtt))
	v.Set("chl", strings.Join(chl, "|"))
	v.Set("chd", "t:"+strings.Join(chd, ","))
	chartURL = "http://chart.oa.com/chart.php?" + v.Encode()
	return
}
