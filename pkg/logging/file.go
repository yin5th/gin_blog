package logging

import (
	"fmt"
	"gin_blog/pkg/setting"
	"time"
)

var (
	LogSavePath = setting.AppSetting.LogSavePath
	LogSaveName = setting.AppSetting.LogSaveName
	LogFileExt  = setting.AppSetting.LogFileExt
	TimeFormat  = setting.AppSetting.TimeFormat
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}
