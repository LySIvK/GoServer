package main

import (
	"fmt"
	"gameserver/modulemgr"
	"gameserver/playermgr"
	"gopath/code.google.com/p/go.net/websocket"
	"loger"
)

func Test() {
	fmt.Println("Hello")
}

type GameServer struct {
	serverID            int                    //! 游戏服务器ID
	serverLimit         int                    //! 游戏服务器限制人数
	addPlayerChannel    chan *PlayerMgr.Player //! 客户端登入通道
	removePlayerChannel chan *PlayerMgr.Player //! 客户端登出通道

	moduleMgr *ModuleMgr.ModuleMgr //! 模块管理器
	playerMgr *PlayerMgr.PlayerMgr //! 玩家管理器
}

//! 初始化服务器
func (self *GameServer) Init(serverID int, limit int) {
	self.serverID = serverID
	self.serverLimit = limit
	self.addPlayerChannel = make(chan *PlayerMgr.Player)
	self.removePlayerChannel = make(chan *PlayerMgr.Player)
	self.playerMgr = PlayerMgr.NewPlayerMgr(self.serverLimit)
	self.moduleMgr = ModuleMgr.NewModuleMgr()
}

//! 服务器连接回调
func (self *GameServer) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {

		//! 创建一个新玩家
		player := PlayerMgr.NewPlayer(ws)
		player.PlayerID = self.playerMgr.CreateNewPlayerID()

		//! 通知频道
		self.addPlayerChannel <- player

		//! 开启协程
		player.Run()

		//! 玩家离开
		self.removePlayerChannel <- player
	}
	return websocket.Handler(connectHandler)
}

//! 服务器监听
func (self *GameServer) Listen() {
	for {
		select {
		case player := <-self.addPlayerChannel:
			loger.Debug("Player connect: %d", player.PlayerID)
			self.playerMgr.AddPlayer(player)
		case player := <-self.removePlayerChannel:
			loger.Debug("Player disconnect: %d", player.PlayerID)
			self.playerMgr.SubPlayer(player)
		}
	}
}

//! 生成一个新的游戏服务器
func NewGameServer(serverID int, limit int) *GameServer {
	server := new(GameServer)
	server.Init(serverID, limit)
	return server
}
