package login

import (
	"gameserver/playermgr"
)

type LoginMgr struct {
	playerInfo PlayerMgr.PlayerInfo
}

func (self *LoginMgr) registerMsg() {
	//! 玩家登录游戏服务器
	PlayerMgr.G_Dispatch.AddMsgRegistryToMap(new(PlayerLogin))

}

func (self *LoginMgr) Init() {
	//! 注册消息
	self.registerMsg()
}

func NewLoginMgr() *LoginMgr {
	mgr := new(LoginMgr)
	mgr.Init()
	return mgr
}
