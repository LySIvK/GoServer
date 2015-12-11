package module

import ()

type LoginMgr struct {
	moduleMgr *ModuleMgr //! 模块管理器指针
}

func (self *LoginMgr) registerMsg() {
	//! 玩家登录游戏服务器
	G_Dispatch.AddMsgRegistryToMap(new(Msg_PlayerLogin))

}

func (self *LoginMgr) Init(moduleMgr *ModuleMgr) {
	//! 保存模块管理器指针
	self.moduleMgr = moduleMgr

	//! 注册消息
	self.registerMsg()
}

func NewLoginMgr(moduleMgr *ModuleMgr) *LoginMgr {
	loginMgr := new(LoginMgr)
	loginMgr.Init(moduleMgr)
	return loginMgr
}
