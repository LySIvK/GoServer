package ModuleMgr

import (
	"gameserver/login"
)

type ModuleMgr struct {
	//! 登录管理器
	loginMgr *login.LoginMgr
}

func (self *ModuleMgr) Init() {
	self.loginMgr = login.NewLoginMgr()
}

func NewModuleMgr() *ModuleMgr {
	mgr := new(ModuleMgr)
	mgr.Init()
	return mgr
}
