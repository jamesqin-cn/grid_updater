package util

import (
	"log"
	"testing"
)

func TestBuildLcImgUrl(t *testing.T) {
	chs := "300x200"
	chtt := "中文标题"
	chmixer := "company"
	chm := "N"
	chdl := []string{"sz", "bj"}
	chxl := []string{"05-04", "05-05", "05-06", "05-07", "05-08"}
	chd := [][]float64{{200, 50, 60, 80, 40}, {50, 60, 100, 40, 20}}
	expectedUrl := "http://chart.oa.com/chart.php?chd=t%3A200.00%2C50.00%2C60.00%2C80.00%2C40.00%7C50.00%2C60.00%2C100.00%2C40.00%2C20.00&chdl=sz%7Cbj&chm=N&chmixer=company&chs=300x200&cht=lc&chtt=%25E4%25B8%25AD%25E6%2596%2587%25E6%25A0%2587%25E9%25A2%2598&chxl=0%3A%7C05-04%7C05-05%7C05-06%7C05-07%7C05-08"

	url, err := BuildLcImgUrl(chs, chtt, chmixer, chm, chdl, chxl, chd)
	log.Println("BuildLcImgUrl(), get url:" + url)
	if err != nil || expectedUrl != url {
		t.Errorf("chart url build failed:%s, want:%s, get:%s", err, expectedUrl, url)
	}
}
