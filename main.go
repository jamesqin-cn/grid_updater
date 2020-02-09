package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"git.oa.com/data_warehouse/util"
)

var (
	mysqlHost     = flag.String("host", "127.0.0.1", "mysql host")
	mysqlPort     = flag.Int("port", 3306, "mysql port")
	mysqlUsername = flag.String("username", "root", "mysql username")
	mysqlPassword = flag.String("password", "123456", "mysql password")
	charset       = flag.String("charset", "utf8mb4", "default charset")

	dataFile = flag.String("datafile", "", "data file which will be load to mysql")
	fieldSep = flag.String("fieldsep", "\t", "field seperator chat of each records")

	showVer = flag.Bool("version", false, "print current verion")
	debug   = flag.Bool("debug", false, "print debug message")
)

func init() {
	flag.Parse()
	if *showVer {
		fmt.Println(util.GetAppInfo())
		os.Exit(0)
	}

	if len(*dataFile) == 0 {
		log.Fatalln("Missing parameter <datafile>")
		return
	}
}

func main() {
	basefile, err := util.NewBaseFile(*dataFile, *fieldSep)
	if err != nil {
		log.Fatalln("Create basefile object failed, err = ", err)
		return
	}
	if basefile.IsEmptyFile() {
		log.Fatalln("Basefile is an empty file")
		return
	}

	mysql := util.NewMySQL(*mysqlUsername, *mysqlPassword, *mysqlHost, *mysqlPort, "", *debug)
	updater, err := NewUpdater(basefile, mysql, *charset)
	if err != nil {
		log.Fatalln("Init updater failed, err = ", err)
		return
	}

	err = updater.Doit()
	if err != nil {
		log.Fatalln("GridUpdater process failed, err = ", err)
		return
	}
}
