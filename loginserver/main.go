package main

import (
	"db"
	"fmt"
	"loger"
	"serverconfig"
	"tool"
)

func main() {

	//! 初始化日志管理器
	loger.InitLoger("./log", loger.LogDebug, true, "loginserver")

	//! 初始化配置文件
	serverconfig.Init()

	//! 初始化数据库管理器
	database_url := fmt.Sprintf("%s:%d", serverconfig.G_Config.Database_IP, serverconfig.G_Config.Database_Port)
	database_PoolLimit := 1024
	db.Init(database_url, database_PoolLimit)

	//! 初始化工具类
	tool.Init()

	//! 初始化服务器
	loginServer := new(LoginServer)
	loginServer.Init(serverconfig.G_Config.LoginServerID)

	//! 开启控制台窗口, 接收用户指令
	tool.StartConsole()

	//! 注册消息并开启服务器
	loger.Info("---------Login Server---------")
	loginServer.RegHttpMsgHandler()
	loginServer.StartServer()

}
