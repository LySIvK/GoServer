package serverconfig

import (
	"encoding/json"
	"io/ioutil"
	"loger"
	"os"
)

type Config struct {
	Database_IP      string `json:"dbip"`              //! 数据库IP
	Database_Port    int    `json:"dbport"`            //! 数据库端口
	LoginServer_IP   string `json:"loginserver_ip"`    //! 登录服务器IP
	LoginServer_Port int    `json:"loginserver_port"`  //! 登录服务器端口
	LoginServerID    int    `json:"loginserver_id"`    //! 登录服务器ID
	LoginServerLimit int    `json:"loginserver_limit"` //! 登陆服务器人数限制
}

//! 全局变量 服务器配置文件
var G_Config Config

const ConfigPath = "./config.json"

//! 载入配置文件
func Init() {
	f, err := os.Open(ConfigPath)
	if err != nil {
		//! 检测配置文件是否存在
		isExist := os.IsExist(err)
		if isExist == true {
			loger.Fatal("Open config file fail. Error: %v", err.Error())
			return
		}

		//! 文件不存在则创建默认配置文件
		CreateDefaultConfigFile()
		f, _ = os.Open(ConfigPath)
	}

	//! 延后关闭文件
	defer f.Close()

	//! 读取配置文件内容
	data, err := ioutil.ReadAll(f)
	if err != nil {
		loger.Error("Read config file fail. Error: %v", err.Error())
		return
	}

	//! 解析Json格式配置
	err = json.Unmarshal(data, &G_Config)
	if err != nil {
		loger.Error("Unmarshal config file fail. Error: %v", err.Error())
		return
	}

	loger.Debug("Load config success: \r\n %s", string(data))
}

//! 创建默认配置文件
func CreateDefaultConfigFile() {
	f, err := os.Create(ConfigPath)
	if err != nil {
		loger.Fatal("Create default config file fail. Error: %v", err.Error())
	}

	//! 设置默认配置参数
	c := Config{
		Database_IP:      "127.0.0.1",
		Database_Port:    27017,
		LoginServer_IP:   "127.0.0.1",
		LoginServer_Port: 9016,
		LoginServerID:    1,
		LoginServerLimit: 5000,
	}

	//! 将配置解析为Json格式
	j, err := json.Marshal(c)
	if err != nil {
		loger.Error("Marshal json fail. Error: %v", err.Error())
	}

	write_num, err := f.Write(j)
	if err != nil {
		loger.Error("Write default config file fail, write %d bytes. Error: %v", write_num, err.Error())
	}

	f.Close()
}
