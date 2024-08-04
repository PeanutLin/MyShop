package common

import "os"

// 判断文件是否存在
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
