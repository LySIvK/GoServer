package main

import (
	"db"
	"flag"
	"fmt"
	"loger"
	"net/http"
	"serverconfig"
	"strconv"
)

func main() {
	//! 初始化日志管理器
	loger.InitLoger("./log", loger.LogDebug, true, "gameserver")

	//! 初始化服务器配置
	serverconfig.Init()

	//! 初始化数据库管理器
	database_url := fmt.Sprintf("%s:%d", serverconfig.G_Config.Database_IP, serverconfig.G_Config.Database_Port)
	database_PoolLimit := 1024
	db.Init(database_url, database_PoolLimit)

	//! 获取服务器地址信息
	var addr = flag.String(serverconfig.G_Config.GameServer_IP, ":"+strconv.Itoa(serverconfig.G_Config.GameServer_Port), "websocket game service address")
	flag.Parse()

	//! 创建游戏服务器
	serverSingleton := NewGameServer(serverconfig.G_Config.GameServerID, serverconfig.G_Config.GameServerLimit)

	//! 设置websocket回调
	onConnectHandler := serverSingleton.GetConnectHandler()
	http.Handle("/", onConnectHandler)

	//! 开启监听
	go serverSingleton.Listen()

	loger.Info("---------Game Server---------")
	loger.Info("---------Addr: %s ---------", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		loger.Fatal("ListenAndServe: ", err)
	}
}
