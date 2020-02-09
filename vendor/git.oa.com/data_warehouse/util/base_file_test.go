package util

import (
	"fmt"
	"testing"
)

func TestParseBaseFileName(t *testing.T) {
	fileName := "/data/a/b/c/test#x_user_summary#20171214#1#0#fund_withdraw.txt"
	base, err := NewBaseFile(fileName, "\t")
	if err != nil {
		t.Errorf("parse base file faild")
	} else {
		if base.DbName() != "test" {
			t.Errorf("Parse DbName failed, want %s, but get %s", "test", base.DbName())
		}
		if base.TableName() != "x_user_summary" {
			t.Errorf("Parse TbName failed, want %s, but get %s", "x_user_summary", base.TableName())
		}
		if base.StatTime() != "2017-12-14" {
			t.Errorf("Parse statTime failed, want %s, but get %s", "2017-12-14", base.StatTime())
		}
		if base.PrimaryKeyNum() != 1 {
			t.Errorf("Parse statTime failed, want %d, but get %d", 1, base.PrimaryKeyNum())
		}
		if base.EraseColumnBeforeUpdate() != false {
			t.Errorf("Parse eraseColumnBeforeUpdate failed, want %t, but get %t", false, base.EraseColumnBeforeUpdate())
		}
		if base.FileDesc() != "fund_withdraw.txt" {
			t.Errorf("Parse eraseColumnBeforeUpdate failed, want %s, but get %s", "fund_withdraw.txt", base.FileDesc())
		}
	}
}

func TestSettings(t *testing.T) {
	fileName := "../grid_updater/testdata/test#x_user_summary#20171214#1#0#fund_withdraw.txt"
	base, _ := NewBaseFile(fileName, "\t")
	if err := base.Open(); err != nil {
		t.Errorf("open file failed.")
	} else {
		settings := base.Settings()
		fmt.Println("settings = ", settings)
		fmt.Println("before_command = ", base.BeforeCommand())
		fmt.Println("after_command = ", base.AfterCommand())

	}
}

func TestColumns(t *testing.T) {
	fileName := "../grid_updater/testdata/test#x_user_summary#20171214#1#0#fund_withdraw.txt"
	base, _ := NewBaseFile(fileName, "\t")
	if err := base.Open(); err != nil {
		t.Errorf("open file failed.")
	} else {
		colsName, colsType, idxs, _ := base.Columns()
		fmt.Printf("colsName = %v, colsType = %v,indexs = %v\n", colsName, colsType, idxs)
	}
}

func TestIsPrimaryKey(t *testing.T) {
	fileName := "../grid_updater/testdata/test#x_user_summary#20171214#1#0#fund_withdraw.txt"
	base, _ := NewBaseFile(fileName, "\t")
	if err := base.Open(); err != nil {
		t.Errorf("open file failed.")
	} else {
		if base.IsPrimaryKey("x_username") == false {
			t.Errorf("TestIsPrimaryKey failed, x_username is primery key, want true, but get false")
		}

		if base.IsPrimaryKey("nokey") == true {
			t.Errorf("TestIsPrimaryKey failed, nokey is primery key, want false, but get true")
		}
	}
}

func TestNextRow(t *testing.T) {
	fileName := "../grid_updater/testdata/test#x_user_summary#20171214#1#0#fund_withdraw.txt"
	base, _ := NewBaseFile(fileName, "\t")
	if err := base.Open(); err != nil {
		t.Errorf("open file failed.")
	} else {
		for {
			row, err := base.NextRow()
			if err != nil || row == nil {
				break
			}
			fmt.Println(row)
		}
	}
}
