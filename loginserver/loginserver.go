package main

import (
	account "loginserver/accountmgr"
	game "loginserver/gameservermgr"
	"net/http"
	"serverconfig"
	"strconv"
	"tool"
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
	http.HandleFunc("/", self.Test)

	//! 玩家登陆
	http.HandleFunc("/login", self.accountMgr.Handler_UserLogin)

	//! 玩家注册
	http.HandleFunc("/register", self.accountMgr.Handler_UserRegister)

	//! 玩家请求服务器列表
	http.HandleFunc("/serverlist", self.accountMgr.Handler_ServerList)

	//! 验证登录
	http.HandleFunc("/verifyuserlogin", self.accountMgr.Handler_VerifyUserLogin)

	//! 注册服务器
	http.HandleFunc("/reggameserver", self.gameServerMgr.Handler_RegisterGameSvr)
}

func (self *LoginServer) StartServer() {
	tool.HttpLimitListen(":"+strconv.Itoa(serverconfig.G_Config.LoginServer_Port), serverconfig.G_Config.LoginServerLimit)
}
