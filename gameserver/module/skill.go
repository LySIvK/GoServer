package module

import (
//"loger"
)

type SkillInfo struct {
	ID   int    `bson:"_id"`
	Name string `json:"name"`
}

type SkillInfoLst map[int]*SkillInfo
type SkillMgr struct {
	moduleMgr *ModuleMgr

	skillLst SkillInfoLst //! 玩家技能信息
}

//! 初始化
func (self *SkillMgr) Init(moduleMgr *ModuleMgr) {
	self.moduleMgr = moduleMgr
	self.skillLst = make(SkillInfoLst)

}

func NewSkillMgr(moduleMgr *ModuleMgr) *SkillMgr {
	mgr := new(SkillMgr)
	mgr.Init(moduleMgr)
	return mgr
}
