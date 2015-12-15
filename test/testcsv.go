package main

import (
	"encoding/csv"
	"io"
	"loger"
	"os"
	"reflect"
	"strconv"
)

type StaticData interface {
	GetPathName() string //! 获取路径名
	GetName() string     //！获取命名
}

type StaticDataMgr struct {
	csvLst map[string]StaticData
}

func (self *StaticDataMgr) Init() {
	self.csvLst = make(map[string]StaticData)
}

func (self *StaticDataMgr) Add(data StaticData) {
	self.csvLst[data.GetName()] = data
}

func CloneType(obj interface{}) interface{} {
	newObj := reflect.New(reflect.TypeOf(obj).Elem()).Elem()
	return newObj.Addr().Interface()
}

func (self *StaticDataMgr) Parse() {
	for _, v := range self.csvLst {
		file, err := os.Open(v.GetPathName())
		if err != nil {
			loger.Fatal("open csv file failed. Error: %s", err.Error())
		}

		reader := csv.NewReader(file)
		lineNumber := 0
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}

				loger.Fatal("read csv error. Error: %s", err.Error())
			}

			if lineNumber > 0 {
				//! 创建一个新的无值结构体
				newData := CloneType(v)

				//! 使用反射机制来进行接口赋值
				ref := reflect.ValueOf(newData).Elem()

				for i := 0; i < ref.NumField(); i++ {
					value := ref.Field(i)
					valueType := value.Kind()
					if valueType == reflect.Int {
						data, _ := strconv.ParseInt(record[i], 10, 64)
						value.SetInt(data)
					} else if valueType == reflect.String {
						data := record[i]
						value.SetString(data)
					}
				}

				loger.Info("%v", newData)

			}
			lineNumber++
		}

		file.Close()
	}
}

func TestLoadCsv() {
	//! 打开csv文件
	file, err := os.Open("test.csv")
	if err != nil {
		loger.Error("open csv file failed. Error: %s", err.Error())
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			loger.Error("read csv error. Error: %s", err.Error())
			return
		}

		loger.Info("record: %v", record)
	}
}

func TestParseCsv() {
	type TestData struct {
		ID    int64
		Name  string
		Level int
		Money int
	}

	file, err := os.Open("test.csv")
	if err != nil {
		loger.Error("open csv file failed. Error: %s", err.Error())
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	lineNumber := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			loger.Error("read csv error. Error: %s", err.Error())
			return
		}

		if lineNumber > 0 {
			//! 遍历字符串slice，取出数据

			var data TestData
			data.ID, _ = strconv.ParseInt(record[0], 10, 64)
			data.Name = record[1]
			data.Level, _ = strconv.Atoi(record[2])
			data.Money, _ = strconv.Atoi(record[3])

			loger.Info("%v", data)
		}

		lineNumber++
	}
}
