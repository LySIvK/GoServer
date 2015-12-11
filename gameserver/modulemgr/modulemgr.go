package ModuleMgr

import (
	"gameserver/login"
	"gameserver/skill"
)

type ModuleMgr struct {
	//! 登录管理器
	loginMgr *login.LoginMgr
	skillMgr *skill.SkillMgr
}

func (self *ModuleMgr) Init() {
	self.loginMgr = login.NewLoginMgr()
	self.skillMgr = skill.NewSkillMgr()
}

func NewModuleMgr() *ModuleMgr {
	mgr := new(ModuleMgr)
	mgr.Init()
	return mgr
}
