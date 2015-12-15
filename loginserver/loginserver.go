package main

import (
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
	account "loginserver/accountmgr"
	game "loginserver/gameservermgr"
	"net/http"
	"serverconfig"
	"strconv"
	//"tool"
)

//! 登录服务器
type LoginServer struct {
	serverID      int                 //! 登陆服务器ID
	gameServerMgr *game.GameServerMgr //! 游戏服务器管理器
	accountMgr    *account.AccountMgr //! 帐号管理器
}

func (self *LoginServer) Init(loginServerID int) {
	//! ID赋值
	self.serverID = loginServerID

	//! 初始化各类管理器
	self.gameServerMgr = game.NewGameServerMgr()
	self.accountMgr = account.NewAccountMgr(self.gameServerMgr)
}

func (self *LoginServer) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (self *LoginServer) RegHttpMsgHandler() {
	//! 测试
	goji.Get("/", self.Test)
	//! http.HandleFunc("/", self.Test)

	//! 玩家登陆
	goji.Post("/login", self.accountMgr.Handler_UserLogin)
	//! http.HandleFunc("/login", self.accountMgr.Handler_UserLogin)

	//! 玩家注册
	goji.Post("/register", self.accountMgr.Handler_UserRegister)
	//! http.HandleFunc("/register", self.accountMgr.Handler_UserRegister)

	//! 玩家请求服务器列表
	goji.Post("/serverlist", self.accountMgr.Handler_ServerList)
	//! http.HandleFunc("/serverlist", self.accountMgr.Handler_ServerList)

	//! 验证登录
	goji.Post("/verifyuserlogin", self.accountMgr.Handler_VerifyUserLogin)
	//! http.HandleFunc("/verifyuserlogin", self.accountMgr.Handler_VerifyUserLogin)

	//! 注册服务器
	goji.Post("/reggameserver", self.gameServerMgr.Handler_RegisterGameSvr)
	//! http.HandleFunc("/reggameserver", self.gameServerMgr.Handler_RegisterGameSvr)
}

func (self *LoginServer) StartServer() {
	//! 使用Goji框架
	listener := bind.Socket(":" + strconv.Itoa(serverconfig.G_Config.LoginServer_Port))
	goji.ServeListener(listener)

	//! 使用原生库HTTP框架
	//tool.HttpLimitListen(":"+strconv.Itoa(serverconfig.G_Config.LoginServer_Port), serverconfig.G_Config.LoginServerLimit)
}
