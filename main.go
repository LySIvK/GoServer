package main

import (
	"loger"
	"serverconfig"
)

func main() {
	//! 初始化日志管理器
	loger.InitLoger("./log", loger.LogDebug, true, "test")

	//! 初始化服务器配置
	serverconfig.Init()

	//! 获取服务器地址信息
	var addr = flag.String("0.0.0.0", ":9000", "http service address")

	//! Websocket
	flag.Parse()

	serverSingleton := NewGameServer(*addr)
	serverSingleton.InitServer()
	onConnectHandler := serverSingleton.GetConnectHandler()
	http.Handle("/", onConnectHandler)
	go serverSingleton.Listen()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		loger.Fatal("ListenAndServe: ", err)
	}
}
