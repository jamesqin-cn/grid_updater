package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"git.oa.com/data_warehouse/util"
)

type Updater struct {
	basefile *util.BaseFile
	mysql    *util.MySQL
	charset  string
    interval int
}

func NewUpdater(reader *util.BaseFile, writer *util.MySQL, charset string, interval int) (updater *Updater, err error) {
	if err = reader.Open(); err != nil {
		return
	}
	if err = writer.Connect(); err != nil {
		return
	}

	return &Updater{
		basefile: reader,
		mysql:    writer,
		charset:  charset,
		interval: interval,
	}, nil
}

func (h *Updater) ExecBeforeCommand() (err error) {
	if len(h.basefile.BeforeCommand()) > 0 {
		_, err = h.mysql.Exec(h.basefile.BeforeCommand())
	}
	return
}

func (h *Updater) ExecAfterCommand() (err error) {
	if len(h.basefile.AfterCommand()) > 0 {
		_, err = h.mysql.Exec(h.basefile.AfterCommand())
	}
	return
}

// 建库
func (h *Updater) BuildDatabase() (err error) {
	if h.mysql.ExistsDatabase(h.basefile.DbName()) == false {
		sql := fmt.Sprintf("CREATE DATABASE %s", h.basefile.DbName())
		_, err = h.mysql.Exec(sql)
	}
	return
}

// 建表
func (h *Updater) BuildTable(charset string) (err error) {
	colsName, colsType, _, err := h.basefile.Columns()
	if err != nil {
		log.Println("BuildTable failed, err = ", err)
		return
	}

	pks := colsName[:h.basefile.PrimaryKeyNum()]
	pkNameStr := "`" + strings.Join(pks, "`,`") + "`"

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`%s`( \n", h.basefile.DbName(), h.basefile.TableName())
	for i := 0; i < h.basefile.PrimaryKeyNum(); i++ {
		fieldName := colsName[i]
		fieldType := colsType[i]
		sql = sql + fmt.Sprintf("  `%s` %s NOT NULL,\n", fieldName, fieldType)
	}
	sql = sql + "  `_created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"
	sql = sql + "  `_updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,\n"
	sql = sql + "  `_deleted_at` datetime DEFAULT NULL,\n"
	sql = sql + "  PRIMARY KEY(" + pkNameStr + ")\n"
	sql = sql + ") ENGINE=InnoDB DEFAULT CHARSET=" + charset

	_, err = h.mysql.Exec(sql)
	return
}

// 判断字段是否属于数字类型
func (h *Updater) IsNumberColumnType(fieldType string) bool {
	fieldType = strings.ToUpper(fieldType)
	numberTypes := []string{"INT", "DECIMAL", "FLOAT", "DOUBLE", "REAL"}
	for _, t := range numberTypes {
		if strings.Contains(fieldType, t) {
			return true
		}
	}

	return false
}

// 创建字段
func (h *Updater) BuildColumns() (err error) {
	colsName, colsType, _, err := h.basefile.Columns()
	if err != nil {
		log.Println("BuildColumns failed, err = ", err)
		return
	}

	sql := fmt.Sprintf("ALTER TABLE `%s`.`%s`", h.basefile.DbName(), h.basefile.TableName())
	hasChange := false
	for i := h.basefile.PrimaryKeyNum(); i < len(colsName); i++ {
		if h.mysql.ExistsColumn(h.basefile.DbName(), h.basefile.TableName(), colsName[i]) == false {
			if hasChange == true {
				sql += ","
			}

			if h.IsNumberColumnType(colsType[i]) {
				//sql += fmt.Sprintf(" ADD COLUMN `%s` %s NOT NULL DEFAULT 0", colsName[i], colsType[i])
				sql += fmt.Sprintf(" ADD COLUMN `%s` %s DEFAULT NULL", colsName[i], colsType[i])
			} else {
				sql += fmt.Sprintf(" ADD COLUMN `%s` %s DEFAULT NULL", colsName[i], colsType[i])
			}

			hasChange = true
		}
	}

	if hasChange {
		_, err = h.mysql.Exec(sql)
	}

	return
}

// 创建索引
func (h *Updater) BuildIndexs() (err error) {
	_, _, idxs, err := h.basefile.Columns()
	if err != nil {
		log.Println("BuildIndexs failed, err = ", err)
		return
	}

	sql := fmt.Sprintf("ALTER TABLE `%s`.`%s`", h.basefile.DbName(), h.basefile.TableName())
	hasChange := false
	for _, idx := range idxs {
		if h.mysql.ExistsIndex(h.basefile.DbName(), h.basefile.TableName(), idx) == false {
			if hasChange == true {
				sql += ","
			}
			sql += fmt.Sprintf(" ADD INDEX(`%s`)", idx)

			hasChange = true
		}
	}

	if hasChange {
		_, err = h.mysql.Exec(sql)
	}

	return
}

// 构建记录
func (h *Updater) UpdateRecords() (err error) {
	colsName, colsType, _, err := h.basefile.Columns()
	if err != nil {
		log.Println("UpdateRecords.Columns failed, err = ", err)
		return
	}

	if h.basefile.EraseColumnBeforeUpdate() {
		for i, col := range colsName {
			if i < h.basefile.PrimaryKeyNum() {
				continue
			}

			var sql string
			if h.IsNumberColumnType(colsType[i]) {
				sql = fmt.Sprintf("UPDATE `%s`.`%s` SET `%s`=0 WHERE `%s`<>0",
					h.basefile.DbName(), h.basefile.TableName(), col, col)
			} else {
				sql = fmt.Sprintf("UPDATE `%s`.`%s` SET `%s`=NULL WHERE `%s` IS NOT NULL",
					h.basefile.DbName(), h.basefile.TableName(), col, col)
			}

			_, err = h.mysql.Update(sql)
			if err != nil {
				log.Println("UpdateRecords.EraseColumnBeforeUpdate failed, err = ", err)
				return
			}
		}
	}

	for {
		row, _ := h.basefile.NextRow()
		if row == nil {
			continue
		}

		h.updateRecord(row)
        time.Sleep(time.Duration(h.interval) * time.Millisecond)
	}

	return
}

func (h *Updater) updateRecord(row map[string]string) (err error) {
	updateCols := make([]string, 0)
	updateVals := make([]string, 0)
	updatePair := make([]string, 0)
	for k, v := range row {
		updateCols = append(updateCols, "`"+k+"`")
		if strings.ToUpper(v) == "NULL" {
			colType, _ := h.basefile.GetColumnTypeByName(k)
			if h.IsNumberColumnType(colType) {
				updateVals = append(updateVals, "0")
				if h.basefile.IsPrimaryKey(k) == false {
					updatePair = append(updatePair, "`"+k+"`=0")
				}
			} else {
				updateVals = append(updateVals, "NULL")
				if h.basefile.IsPrimaryKey(k) == false {
					updatePair = append(updatePair, "`"+k+"`=NULL")
				}
			}
		} else if strings.ToUpper(v) == "" {
			colType, _ := h.basefile.GetColumnTypeByName(k)
			if h.IsNumberColumnType(colType) {
				updateVals = append(updateVals, "0")
				if h.basefile.IsPrimaryKey(k) == false {
					updatePair = append(updatePair, "`"+k+"`=0")
				}
			} else {
				updateVals = append(updateVals, "''")
				if h.basefile.IsPrimaryKey(k) == false {
					updatePair = append(updatePair, "`"+k+"`=''")
				}
			}
		} else {
			updateVals = append(updateVals, "'"+h.mysql.AddSlashes(v)+"'")
			if h.basefile.IsPrimaryKey(k) == false {
				updatePair = append(updatePair, "`"+k+"`='"+h.mysql.AddSlashes(v)+"'")
			}
		}
	}
	sql := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s,`_updated_at`='%s'",
		h.basefile.DbName(),
		h.basefile.TableName(),
		strings.Join(updateCols, ","),
		strings.Join(updateVals, ","),
		strings.Join(updatePair, ","),
		util.GetNow())
	_, err = h.mysql.Insert(sql)
	if err != nil {
		log.Println("UpdateRecord.Insert failed, err = ", err)
		return
	}
	return
}

func (h *Updater) Doit() (err error) {
	if err = h.BuildDatabase(); err != nil {
		return
	}
	if err = h.mysql.UseDb(h.basefile.DbName()); err != nil {
		return
	}
	if err = h.BuildTable(h.charset); err != nil {
		return
	}
	if err = h.BuildColumns(); err != nil {
		return
	}
	if err = h.BuildIndexs(); err != nil {
		return
	}
	if err = h.ExecBeforeCommand(); err != nil {
		return
	}
	h.UpdateRecords()
	if err = h.ExecAfterCommand(); err != nil {
		return
	}

	return
}

func (h *Updater) Close() {
	h.basefile.Close()
}
