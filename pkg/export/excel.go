package export

import "go-gin-demo/pkg/setting"

const EXT = ".xlsx"

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

//excel路径
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

//excel全路径
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
