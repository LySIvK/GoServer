package main

import (
	"fmt"
)

func Test() {
	fmt.Println("Hello")
}

type GameServer struct {
	serverID            int      //! 游戏服务器ID
	addPlayerChannel    chan int //! 客户端登入通道
	removePlayerChannel chan int //! 客户端登出通道
}
