package util

import (
	"fmt"
	"testing"
)

func TestNewWithConfigSvc(t *testing.T) {
	dbUtil := NewMySQLWithConfigSvc("online_ro", true)
	defer dbUtil.Close()

	rows, _, _ := dbUtil.Query("SHOW DATABASES")

	for lineNo, row := range rows {
		fmt.Printf("Line %d (%d columns)\n", lineNo, len(row))
		for k, v := range row {
			fmt.Printf("    %s = %s\n", k, v)
		}
	}
}

func TestQuery(t *testing.T) {
	dbUtil := NewMySQL("root", "123456", "mysql.proxy.oa.com", 3306, "", true)
	defer dbUtil.Close()

	rows, _, _ := dbUtil.Query("SHOW DATABASES")

	for lineNo, row := range rows {
		fmt.Printf("Line %d (%d columns)\n", lineNo, len(row))
		for k, v := range row {
			fmt.Printf("    %s = %s\n", k, v)
		}
	}
}

func TestExistDatabase(t *testing.T) {
	dbUtil := NewMySQL("root", "123456", "127.0.0.1", 3306, "test", false)
	defer dbUtil.Close()

	if dbUtil.ExistsDatabase("test") == false {
		t.Errorf("ExistsDatabase() failed, database test exist actually")
	}

	if dbUtil.ExistsDatabase("bad_dbname") == true {
		t.Errorf("ExistsDatabase() failed, database bad_dbname not exist actually")
	}
}

func TestExistTable(t *testing.T) {
	dbUtil := NewMySQL("root", "123456", "127.0.0.1", 3306, "test", false)
	defer dbUtil.Close()

	if dbUtil.ExistsTable("test", "log") == false {
		t.Errorf("ExistsTable() failed, table test.log exist actually")
	}

	if dbUtil.ExistsTable("test", "bad_table") == true {
		t.Errorf("ExistsTable() failed, table test.bad_table not exist actually")
	}
}

func TestExistsColumn(t *testing.T) {
	dbUtil := NewMySQL("root", "123456", "127.0.0.1", 3306, "mysql", false)
	defer dbUtil.Close()

	if dbUtil.ExistsColumn("test", "log", "log_id") == false {
		t.Errorf("ExistsColumn() failed, column test.log.log_id exist actually")
	}

	if dbUtil.ExistsColumn("test", "bad_table", "bad_column") == true {
		t.Errorf("ExistsColumn() failed, column test.bad_table.bad_column not exist actually")
	}

	if dbUtil.ExistsColumn("test", "log", "bad_column") == true {
		t.Errorf("ExistsColumn() failed, column test.log.bad_column not exist actually")
	}
}
