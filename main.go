package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//! 删除文件
func WalkDir(path string, suffix string) ([]string, error) {
	files := []string{}
	suffix = strings.ToUpper(suffix)

	err := filepath.Walk(path, //! 遍历目录
		func(path string, info os.FileInfo, err error) error {
			//! 忽略目录
			if info.IsDir() == true {
				return nil
			}

			if strings.HasSuffix(strings.ToUpper(info.Name()), suffix) {
				files = append(files, path)
			}
			return nil
		})

	return files, err
}

//! 删除日志文件
func DeleteLogFile() {

	//! 遍历目录
	files, err := WalkDir("./", ".log")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//! 删除文件
	for _, v := range files {
		err := os.Remove(v)
		if err == nil {
			fmt.Println("Remove file: ", v)
		}
	}
}

//! 关闭服务器
func main() {

}
