package util

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dbHandle  *sql.DB
	dbOptions *DBOptions
	debug     bool
}

type DBOptions struct {
	DataSource   string `json:"dataSource"`
	Driver       string `json:"driver"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DBName       string `json:"dbName"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

func (o *DBOptions) IsInvalid() bool {
	if o.Driver == "" {
		return true
	}
	if o.DataSource == "" && o.Host == "" {
		return true
	}
	return false
}

func NewMySQL(user string, password string, host string, port int, dbname string, debug bool) (db *MySQL) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=false&loc=Local", user, password, host, port, dbname)
	return &MySQL{
		dbHandle: nil,
		dbOptions: &DBOptions{
			Driver:       "mysql",
			DataSource:   dsn,
			MaxIdleConns: 100,
			MaxOpenConns: 300,
		},
		debug: debug,
	}
}

func NewMySQLWithConfigSvc(cfgPath string, debug bool) (db *MySQL) {
	opts, _ := LoadDBConfig(cfgPath)
	opts.DataSource = strings.Replace(opts.DataSource, "parseTime=true", "parseTime=false", -1)
	return &MySQL{
		dbHandle:  nil,
		dbOptions: opts,
		debug:     debug,
	}
}

func (h *MySQL) SetDebug(debug bool) {
	h.debug = debug
}

func (h *MySQL) Connect() (err error) {
	if h.dbHandle != nil {
		return
	}

	// connect to mysql server
	db, err := sql.Open(h.dbOptions.Driver, h.dbOptions.DataSource)
	if err != nil {
		log.Println("Open mysql pool handler failed, err = ", err)
		return err
	}

	// connect poll setting
	db.SetMaxIdleConns(h.dbOptions.MaxIdleConns)
	db.SetMaxOpenConns(h.dbOptions.MaxOpenConns)
	db.SetConnMaxLifetime(3600 * time.Second)

	h.dbHandle = db
	return
}

func (h *MySQL) Ping() (err error) {
	if err = h.Connect(); err != nil {
		return
	}

	err = h.dbHandle.Ping()
	if err != nil {
		log.Println("Failed to ping mysql: %s", err)
		return
	}

	return
}

func (h *MySQL) Query(query string, args ...interface{}) (results []map[string]string, cols []string, err error) {
	if err = h.Connect(); err != nil {
		return
	}

	if h.debug {
		log.Printf("[SQL] query = %s, args = %v\n", query, args)
	}

	// do query
	rows, err := h.dbHandle.Query(query, args...)
	if err != nil {
		log.Printf("Query sql failed, sql = %s, args = %v, err = %v\n", query, args, err)
		return
	}
	defer rows.Close()

	// prepare resultset variables
	cols, err = rows.Columns()
	if err != nil {
		log.Println("Fetch column failed, err = ", err)
		return
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			log.Println("MySQL get next rows failed, err = %v\n", err)
			continue
		}

		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}

	return
}

func (h *MySQL) Insert(query string, args ...interface{}) (lastInsertId int64, err error) {
	if err = h.Connect(); err != nil {
		return
	}

	res, err := h.Exec(query, args...)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func (h *MySQL) Update(query string, args ...interface{}) (rowsAffected int64, err error) {
	if err = h.Connect(); err != nil {
		return
	}

	res, err := h.Exec(query, args...)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (h *MySQL) Delete(query string, args ...interface{}) (rowsAffected int64, err error) {
	if err = h.Connect(); err != nil {
		return
	}

	res, err := h.Exec(query, args...)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (h *MySQL) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	if err = h.Connect(); err != nil {
		return
	}

	if h.debug {
		log.Printf("[SQL] query = %s, args = %v\n", query, args)
	}

	res, err = h.dbHandle.Exec(query, args...)
	if err != nil {
		log.Printf("Execute sql failed, sql = %s, args = %v, err = %v\n", query, args, err)
	}
	return
}

func (h *MySQL) Exists(query string, args ...interface{}) (exists bool) {
	rows, _, err := h.Query(query, args...)
	if err != nil {
		return false
	}
	if len(rows) == 0 {
		return false
	}
	return true
}

func (h *MySQL) UseDb(dbName string) (err error) {
	cfg, err := mysql.ParseDSN(h.dbOptions.DataSource)
	h.dbOptions.DataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=false&loc=Local",
		cfg.User, cfg.Passwd, cfg.Addr, dbName)
	h.dbHandle = nil
	h.Connect()
	return
}

func (h *MySQL) ExistsDatabase(dbName string) (exists bool) {
	return h.Exists("SHOW DATABASES LIKE '" + dbName + "'")
}

func (h *MySQL) ExistsTable(dbName string, tbName string) (exists bool) {
	return h.Exists("SHOW TABLES FROM `" + dbName + "` LIKE '" + tbName + "'")
}

func (h *MySQL) ExistsColumn(dbName string, tbName string, col string) (exists bool) {
	return h.Exists("SHOW COLUMNS FROM `" + dbName + "`.`" + tbName + "` WHERE field = '" + col + "'")
}

func (h *MySQL) ExistsIndex(dbName string, tbName string, idx string) (exists bool) {
	return h.Exists("SHOW INDEX FROM `" + dbName + "`.`" + tbName + "` WHERE Column_name = '" + idx + "'")
}

// 转义：引号、双引号添加反斜杠
func (h *MySQL) AddSlashes(val string) string {
	val = strings.Replace(val, "\"", "\\\"", -1)
	val = strings.Replace(val, "'", "\\'", -1)
	return val
}

// 反转义：引号、双引号去除反斜杠
func (h *MySQL) StripSlashes(val string) string {
	val = strings.Replace(val, "\\\"", "\"", -1)
	val = strings.Replace(val, "\\'", "'", -1)
	return val
}

func (h *MySQL) Close() {
	if h.dbHandle != nil {
		h.dbHandle.Close()
		h.dbHandle = nil
	}
}
