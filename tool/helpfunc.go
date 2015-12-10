package tool

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

//! 获取当前文件执行的路径
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}

//! 检查目录是否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	return true
}

//! 加密MD5
func MD5(code string) string {
	h := md5.New()
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil))
}

//! 克隆对象,只有类型,没有值,值是空
func CloneType(obj interface{}) interface{} {
	newObj := reflect.New(reflect.TypeOf(obj).Elem()).Elem()
	return newObj.Addr().Interface()
}
