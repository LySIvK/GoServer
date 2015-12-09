package main

import (
	"fmt"
	"gopath/code.google.com/p/go.net/websocket"
	"loger"
)

func Test() {
	fmt.Println("Hello")
}

type GameServer struct {
	serverID            int          //! 游戏服务器ID
	serverAddr          string       //! 游戏服务器地址
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
	connectHandler := func(ws *websocket.Conn) {
		player := NewPlayer(ws)

		self.addPlayerChannel <- player
		player.Run()
		self.removePlayerChannel <- player
	}
	return websocket.Handler(connectHandler)
}

//! 服务器监听
func (self *GameServer) Listen() {
	for {
		select {
		case player := <-self.addPlayerChannel:
			loger.Debug("Player connect")
		case player := <-self.removePlayerChannel:
			loger.Debug("Player disconnect")
		}
	}
}

//! 生成一个新的游戏服务器
func NewGameServer(addr string, serverID int) *GameServer {
	server := new(GameServer)
	server.serverID = serverID
	server.serverAddr = addr
	server.addPlayerChannel = make(chan *Player)
	server.removePlayerChannel = make(chan *Player)
	return server
}
