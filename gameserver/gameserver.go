package main

import (
	"fmt"
	"gopath/code.google.com/p/go.net/websocket"
)

func Test() {
	fmt.Println("Hello")
}

type GameServer struct {
	serverID            int          //! 游戏服务器ID
	addPlayerChannel    chan *Player //! 客户端登入通道
	removePlayerChannel chan *Player //! 客户端登出通道
}

//! 初始化服务器
func (self *GameServer) Init(serverID int) {
	self.serverID = serverID
	self.addPlayerChannel = make(chan *Player)
	self.removePlayerChannel = make(chan *Player)
}

//! 服务器连接回调
func (self *GameServer) GetConnectHandler() websocket.Handler {

}

//! 生成一个新的游戏服务器
func NewGameServer() *GameServer {
	server := new(GameServer)
	return server
}
