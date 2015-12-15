package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"gameserver/module"
	"loger"
	"loginserver/msg"
	"net/http"
	"serverconfig"
	"time"
)

func Test() {
	fmt.Println("Hello")
}

type GameServer struct {
	serverID            int                 //! 游戏服务器ID
	serverLimit         int                 //! 游戏服务器限制人数
	addPlayerChannel    chan *module.Player //! 客户端登入通道
	removePlayerChannel chan *module.Player //! 客户端登出通道

	moduleMgr *module.ModuleMgr //! 模块管理器
	playerMgr *module.PlayerMgr //! 玩家管理器
}

//! 初始化服务器
func (self *GameServer) Init(serverID int, limit int) {
	self.serverID = serverID
	self.serverLimit = limit
	self.addPlayerChannel = make(chan *module.Player)
	self.removePlayerChannel = make(chan *module.Player)
	self.playerMgr = module.NewPlayerMgr(self.serverLimit)
	self.moduleMgr = module.NewModuleMgr(self.playerMgr)

	//! 注册游戏服务器
	go self.LoopHeart()
}

//! 三十秒发送一次服务器信息给登陆服务器
func (self *GameServer) LoopHeart() {
	sendTime := time.Tick(30 * time.Second)
	for {
		self.RegisterHeart()
		<-sendTime
	}
}

//! 注册登录服务器
func (self *GameServer) RegisterHeart() {
	//! 拼接登陆服务器地址
	loginServerUrl := fmt.Sprintf("http://%s:%d", serverconfig.G_Config.LoginServer_IP, serverconfig.G_Config.LoginServer_Port)
	loginServerUrl += "/reggameserver"
	var req msg.Msg_RegisterGameServer_Req
	req.ServerID = self.serverID
	req.IsNew = serverconfig.G_Config.GameServer_New
	req.ServerName = serverconfig.G_Config.GameServerName
	req.ServerIP = serverconfig.G_Config.GameServer_IP
	req.PlayerNum = self.playerMgr.GetOnlinePlayerCount()

	//! 转换为json
	b, err := json.Marshal(&req)
	if err != nil {
		loger.Error("Register game server fail. Error: %s", err.Error())
		return
	}

	//! Post消息给登录服务器
	resp, err := http.Post(loginServerUrl, "Text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("Post msg to login server fail. Error: %s", err.Error())
		return
	}

	resp.Body.Close()

}

//! 服务器连接回调
func (self *GameServer) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {

		//! 创建一个新玩家
		player := module.NewPlayer(ws, self.playerMgr)
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
			self.playerMgr.AddPlayerCount(player)
		case player := <-self.removePlayerChannel:
			loger.Debug("Player disconnect: %d", player.PlayerID)
			self.playerMgr.SubPlayerCount(player)
		}
	}
}

//! 生成一个新的游戏服务器
func NewGameServer(serverID int, limit int) *GameServer {
	server := new(GameServer)
	server.Init(serverID, limit)
	return server
}
