package login

import (
	"gameserver/playermgr"
)

var loginMgr *LoginMgr

type LoginMgr struct {
}

func (self *LoginMgr) registerMsg() {
	//! 玩家登录游戏服务器
	PlayerMgr.G_Dispatch.AddMsgRegistryToMap(new(Msg_PlayerLogin))

}

func (self *LoginMgr) Init() {
	//! 注册消息
	self.registerMsg()
}

func NewLoginMgr() *LoginMgr {
	loginMgr = new(LoginMgr)
	loginMgr.Init()
	return loginMgr
}
