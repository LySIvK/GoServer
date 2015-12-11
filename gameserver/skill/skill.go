package skill

import (
	"loger"
)

type SkillMgr struct {
}

func (self *SkillMgr) OnlyTest() {
	loger.Debug("Hello")
}

func (self *SkillMgr) Init() {

}

func NewSkillMgr() *SkillMgr {
	mgr := new(SkillMgr)
	mgr.Init()
	return mgr
}
