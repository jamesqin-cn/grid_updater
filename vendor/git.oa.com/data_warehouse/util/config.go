package util

import (
	"fmt"
	"reflect"
)

//LoadDBConfig 载入DB配置
func LoadDBConfig(name string) (opts *DBOptions, err error) {
	opts = &DBOptions{}

	// 先尝试应用程序的配置目录，http://config.svc/<SvcName>-db_<name>.json
	// 相应git路径是 http://config.svc/config/<SvcName>/<SvcName|application>-db_<name>.<yml|json>
	err = LoadConfig("db_"+name, opts)
	if err != nil || opts.IsInvalid() {
		// 找不到再尝试公共DB配置，http://config.svc/db-<name>.json
		// 相应git路径是 http://config.svc/config/db/db-<name>.<yml|json>
		err = LoadConfig("/db-"+name, opts)
	}
	if err != nil || opts.IsInvalid() {
		err = fmt.Errorf("db:%s config is invalid", name)
		return
	}
	return
}

// LoadConfig 从配置中心加载配置
func LoadConfig(path string, v interface{}) error {
	// 以json格式加载（在git存放的也许是yml格式，config.svc会根据存储后缀解析，根据请求后缀转化
	uri := "http://config.svc" + path + ".json"
	if err := LoadJsonFromURL(uri, &v); err != nil {
		fmt.Errorf("[LoadConfig] uri: %s, err: %s", uri, err)
		return err
	}
	fmt.Println("[LoadConfig] uri: %s, v: %+v, t: %v", uri, v, reflect.TypeOf(v))
	return nil
}
