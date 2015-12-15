package module

import (
	"gameserver/staticdata"
)

type ModuleMgr struct {
	playerMgr     *PlayerMgr                //! 玩家管理器
	loginMgr      *LoginMgr                 //! 登录管理器
	skillMgr      *SkillMgr                 //! 技能管理器
	staticDataMgr *staticdata.StaticDataMgr //! 静态数据管理器
}

func (self *ModuleMgr) Init(playerMgr *PlayerMgr) {
	self.playerMgr = playerMgr
	self.playerMgr.SetModuleMgrPoint(self)
	self.loginMgr = NewLoginMgr(self)
	self.skillMgr = NewSkillMgr(self)
	self.staticDataMgr = staticdata.NewStaticDataMgr()
}

func NewModuleMgr(playerMgr *PlayerMgr) *ModuleMgr {
	mgr := new(ModuleMgr)
	mgr.Init(playerMgr)
	return mgr
}
