package util

import "encoding/json"

var (
	APP_COMMIT_ID  string
	APP_BUILD_DATE string
)

func GetAppInfo() string {
	appInfo := map[string]string{
		"build_date": APP_BUILD_DATE,
		"commit_id":  APP_COMMIT_ID,
	}
	info, _ := json.Marshal(appInfo)
	return string(info)
}
