package util

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type BaseFile struct {
	fileHandle              *os.File
	scannerHandle           *bufio.Scanner
	fileName                string
	fieldSep                string
	settings                map[string]string
	colsName                []string
	colsType                []string
	idxsName                []string
	dbName                  string
	tbName                  string
	statTime                string
	primaryKeyNum           int
	eraseColumnBeforeUpdate bool
	fileDesc                string
	isEmptyFile             bool
}

const (
	BASE_FILE_SETTINGS_BEFORE_COMMAND = "before_command"
	BASE_FILE_SETTINGS_AFTER_COMMAND  = "after_command"
	BASE_FILE_INDEX_COLUMN_FLAG       = "!"
)

func NewBaseFile(fileName string, fieldSep string) (basefile *BaseFile, err error) {
	dbName, tbName, statTime, primaryKeyNum, eraseColumnBeforeUpdate, desc, err := ParseBaseFileName(fileName)
	if err != nil {
		return
	}

	return &BaseFile{
		fileHandle:              nil,
		scannerHandle:           nil,
		fileName:                fileName,
		fieldSep:                fieldSep,
		dbName:                  dbName,
		tbName:                  tbName,
		settings:                make(map[string]string),
		colsName:                nil,
		colsType:                nil,
		idxsName:                nil,
		statTime:                statTime,
		primaryKeyNum:           primaryKeyNum,
		eraseColumnBeforeUpdate: eraseColumnBeforeUpdate,
		fileDesc:                desc,
		isEmptyFile:             false,
	}, nil
}

func ParseBaseFileName(fileName string) (dbName, tbName, statTime string, primaryKeyNum int, eraseColumnBeforeUpdate bool, desc string, err error) {
	f := func(c rune) bool {
		return c == '#'
	}
	baseName := filepath.Base(fileName)
	cols := strings.FieldsFunc(baseName, f)
	if len(cols) < 5 {
		err = errors.New("Base file name is wrong format, must be <db>#<table>#<date>#<keyNum>#<eraseColumn>#<desc>.<ext>")
		return
	}

	dbName = cols[0]
	tbName = cols[1]

	if len(cols[2]) == 8 {
		statTime = GetFormatedDateTime(cols[2], FORMAT_FULL_DATE)
	} else {
		statTime = GetFormatedDateTime(cols[2], FORMAT_FULL_DATETIME)
	}

	primaryKeyNum, _ = strconv.Atoi(cols[3])

	if cols[4] == "0" {
		eraseColumnBeforeUpdate = false
	} else {
		eraseColumnBeforeUpdate = true
	}

	if len(cols) == 6 {
		desc = cols[5]
	} else {
		desc = ""
	}

	return
}

func (h *BaseFile) ParseSettingLine(line string) (err error) {
	kvPair := strings.SplitN(line, "=", 2)
	if len(kvPair) != 2 {
		err = errors.New("Base file setting line syntax error, <key>=<val> expected")
	} else {
		k := strings.ToLower(strings.TrimSpace(kvPair[0]))
		v := strings.TrimSpace(kvPair[1])
		h.settings[k] = v
	}
	return
}

func (h *BaseFile) ParseColumn(mixColName string) (colName string, colType string, needIndex bool, err error) {
	arr := strings.Split(mixColName, ":")
	switch len(arr) {
	case 1:
		colName = mixColName
		colType = strings.ToUpper("VARCHAR(128)")
	case 2:
		colName = arr[0]
		colType = strings.ToUpper(arr[1])
	default:
		err = errors.New("Base file column syntax error, column name is" + mixColName)
	}

	if strings.HasPrefix(colName, "!") {
		needIndex = true
		colName = colName[1:]
	} else {
		needIndex = false
	}
	return
}

func (h *BaseFile) Open() (err error) {
	file, err := os.Open(h.fileName)
	if err != nil {
		log.Println("Can't not open file ", h.fileName, ", err = ", err)
		return
	}

	h.fileHandle = file
	h.scannerHandle = bufio.NewScanner(file)

	for {
		if h.scannerHandle.Scan() {
			line := strings.TrimSpace(h.scannerHandle.Text())

			// 碰到注释行，则视为预置指令按KV格式进行解析
			if strings.HasPrefix(line, "#") {
				h.ParseSettingLine(line[1:])
				continue
			}
			if strings.HasPrefix(line, "//") {
				h.ParseSettingLine(line[2:])
				continue
			}

			// 按字段列进行解析
			colsMix := strings.Split(line, h.fieldSep)
			h.colsName = make([]string, 0)
			h.colsType = make([]string, 0)
			h.idxsName = make([]string, 0)

			for _, item := range colsMix {
				colName, colType, needIndex, err := h.ParseColumn(item)
				if err == nil {
					h.colsName = append(h.colsName, colName)
					h.colsType = append(h.colsType, colType)
					if needIndex {
						h.idxsName = append(h.idxsName, colName)
					}
				} else {
					err = errors.New("Base file column syntax error, column name is" + item)
				}
			}
			return
		} else {
			h.isEmptyFile = true
			err = errors.New("Missing column line in the file " + h.fileName)
			break
		}
	}

	return
}

func (h *BaseFile) BeforeCommand() string {
	if v, ok := h.settings[BASE_FILE_SETTINGS_BEFORE_COMMAND]; ok {
		return v
	}
	return ""
}

func (h *BaseFile) AfterCommand() string {
	if v, ok := h.settings[BASE_FILE_SETTINGS_AFTER_COMMAND]; ok {
		return v
	}
	return ""
}

func (h *BaseFile) IsEmptyFile() bool {
	return h.isEmptyFile
}

func (h *BaseFile) Settings() map[string]string {
	return h.settings
}

func (h *BaseFile) DbName() string {
	return h.dbName
}

func (h *BaseFile) TableName() string {
	return h.tbName
}

func (h *BaseFile) StatTime() string {
	return h.statTime
}

func (h *BaseFile) PrimaryKeyNum() int {
	return h.primaryKeyNum
}

func (h *BaseFile) GetColumnIdxByName(findme string) (idx int, err error) {
	for i, colName := range h.colsName {
		if colName == findme {
			return i, nil
		}
	}
	idx = -1
	err = errors.New("Column name " + findme + " not found")
	return
}

func (h *BaseFile) GetColumnTypeByName(findme string) (colType string, err error) {
	idx, err := h.GetColumnIdxByName(findme)
	if err != nil {
		return
	}
	return h.colsType[idx], nil
}

func (h *BaseFile) IsPrimaryKey(findme string) bool {
	idx, err := h.GetColumnIdxByName(findme)
	if err == nil && idx < h.primaryKeyNum {
		return true
	}
	return false
}

func (h *BaseFile) EraseColumnBeforeUpdate() bool {
	return h.eraseColumnBeforeUpdate
}

func (h *BaseFile) FileDesc() string {
	return h.fileDesc
}

func (h *BaseFile) Columns() (cols []string, types []string, idxs []string, err error) {
	if h.colsName != nil {
		return h.colsName, h.colsType, h.idxsName, nil
	}

	return nil, nil, nil, errors.New("missing columns")
}

func (h *BaseFile) NextRow() (row map[string]string, err error) {
	if h.scannerHandle.Scan() {
		cols := strings.Split(h.scannerHandle.Text(), h.fieldSep)
		if len(cols) != len(h.colsName) {
			err = errors.New("Column number not match with first line, want " + strconv.Itoa(len(h.colsName)) + " , get " + strconv.Itoa(len(cols)))
			return
		}
		row = make(map[string]string)
		for idx, data := range cols {
			k := h.colsName[idx]
			row[k] = data
		}
		return
	}

	// 异常处理
	if err = h.scannerHandle.Err(); err != nil {
		log.Println("Reading file failed, err = ", err)
        return
	}
   
    // end of file
	row = nil
	err = nil
	return
}

func (h *BaseFile) Close() {
	if h.fileHandle != nil {
		h.fileHandle.Close()
		h.fileHandle = nil
	}
}
