package staticdata

import (
	"encoding/csv"
	"io"
	"loger"
	"os"
	"reflect"
	"strconv"
	"tool"
)

type StaticData interface {
	GetFilePath() string //! 获取路径
	GetName() string     //! 获取静态数据名
}

type DataMap map[int]interface{}
type DataMapMgr map[string]DataMap

//! 静态数据管理器
type StaticDataMgr struct {
	csvLst        []StaticData //! 加载静态表
	staticDataMap DataMapMgr   //! 静态数据储存
}

func (self *StaticDataMgr) Init() {
	self.staticDataMap = make(DataMapMgr)
	self.initCsvFile()
}

//！添加静态表数据
func (self *StaticDataMgr) initCsvFile() {
	mall := new(Mall_Data)
	self.csvLst = append(self.csvLst, mall)

	vip := new(Vip_Data)
	self.csvLst = append(self.csvLst, vip)
}

//! 解析CSV
func (self *StaticDataMgr) ParseCsv() {
	for _, v := range self.csvLst {
		//! 打开文件
		file, err := os.Open(v.GetFilePath())
		if err != nil {
			loger.Fatal("Open csv file failed. Error: %s  File: %s", err.Error(), v.GetFilePath())
		}

		//! 解析文件
		reader := csv.NewReader(file)
		lineNumber := 0
		for {
			record, err := reader.Read()
			if err != nil {
				//! 判断是否读取完毕
				if err == io.EOF {
					break
				}
				loger.Fatal("Read file failed. Eror: %s  File: %s", err.Error(), v.GetFilePath())
			}

			//! 跳过第一行
			if lineNumber > 0 {
				//! 创建一个新的无值结构体
				newData := tool.CloneType(v)

				//! 使用反射机制来进行接口赋值
				ref := reflect.ValueOf(newData).Elem()

				index := 0
				for i := 0; i < ref.NumField(); i++ {
					//! 获取接口类型与值
					value := ref.Field(i)
					valueType := value.Kind()

					//	loger.Info("type: %v   value: %s", valueType, record[i])

					//! 判断类型
					switch valueType {
					case reflect.Int: //! 整型
						data, err := strconv.ParseInt(record[i], 10, 0)
						if err != nil {
							loger.Fatal("Strconv failed. Parse string: %s error: %s", record[i], err.Error())
						}
						if i == 0 {
							//! 注意!! 默认CSV文件首列必须为ID索引
							index = int(data)
						}

						value.SetInt(data)

					case reflect.String: //! 字符串
						value.SetString(record[i])

					case reflect.Float32: //! 浮点型
						data, err := strconv.ParseFloat(record[i], 32)
						if err != nil {
							loger.Fatal("Strconv failed. Parse string: %s error: %s", record[i], err.Error())
						}
						value.SetFloat(data)

					case reflect.Float64: //! 浮点型64位
						data, err := strconv.ParseFloat(record[i], 64)
						if err != nil {
							loger.Fatal("Strconv failed. Parse string: %s error %s", record[i], err.Error())
						}
						value.SetFloat(data)
					}
				}

				//! 解析完毕，放入Map
				if self.staticDataMap[v.GetName()] == nil {
					self.staticDataMap[v.GetName()] = make(DataMap)
				}

				dataMap := self.staticDataMap[v.GetName()]
				dataMap[index] = newData
				//loger.Info("Index: %d  data: %v", index, newData)
			}

			lineNumber++
		}
	}
}

//! 获取一条数据
func (self *StaticDataMgr) GetStaticData(name string, index int) interface{} {
	dataMap, exist := self.staticDataMap[name]
	if exist == false {
		return nil
	}

	info, ok := dataMap[index]
	if ok == false {
		return nil
	}
	return info
}

//! 获取所有数据
func (self *StaticDataMgr) GetStaticDataMap(name string) *DataMap {
	dataMap, exist := self.staticDataMap[name]
	if exist == false {
		return nil
	}
	return &dataMap
}

//! 获取商城表
func (self *StaticDataMgr) GetMallDataIndex(index int) *Mall_Data {
	data := self.GetStaticData("mall", index)
	if data == nil {
		return nil
	}

	return data.(*Mall_Data)
}

//! 获取商城表所有数据
func (self *StaticDataMgr) GetMallDataAll() []*Mall_Data {
	dataMap := self.GetStaticDataMap("mall")

	dataLst := []*Mall_Data{}
	for _, v := range *dataMap {
		dataLst = append(dataLst, v.(*Mall_Data))
	}

	return dataLst
}

//! 获取VIP信息表
func (self *StaticDataMgr) GetVipDataIndex(index int) *Vip_Data {
	data := self.GetStaticData("vip", index)
	if data == nil {
		return nil
	}

	return data.(*Vip_Data)
}

//! 获取VIP信息表所有数据
func (self *StaticDataMgr) GetVipDataAll() []*Vip_Data {
	dataMap := self.GetStaticDataMap("vip")

	dataLst := []*Vip_Data{}
	for _, v := range *dataMap {
		dataLst = append(dataLst, v.(*Vip_Data))
	}

	return dataLst
}

//! 生成新的静态数据管理器
func NewStaticDataMgr() *StaticDataMgr {
	dataMgr := new(StaticDataMgr)
	dataMgr.Init()
	return dataMgr
}
