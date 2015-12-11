package module

import ()

type ModuleMgr struct {
	playerMgr *PlayerMgr
	loginMgr  *LoginMgr
	skillMgr  *SkillMgr
}

func (self *ModuleMgr) Init(playerMgr *PlayerMgr) {
	self.playerMgr = playerMgr
	self.loginMgr = NewLoginMgr(self)
	self.skillMgr = NewSkillMgr(self)
}

func NewModuleMgr(playerMgr *PlayerMgr) *ModuleMgr {
	mgr := new(ModuleMgr)
	mgr.Init(playerMgr)
	return mgr
}
