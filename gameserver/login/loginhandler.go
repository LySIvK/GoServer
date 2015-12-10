package login

import (
	"gameserver/playermgr"
	"loger"
)

type PlayerLogin struct {
	PlayerMgr.MsgHead `json:"head"`
	AccountID         int64
	LoginKey          string
}

func (self *PlayerLogin) GetTypeAndAction() (string, string) {
	return "login", "playerlogin"
}

func (self *PlayerLogin) checkAction(player *PlayerMgr.Player) bool {
	loger.Debug("Hey,")
	return true
}

func (self *PlayerLogin) payAction(player *PlayerMgr.Player) bool {
	loger.Debug("Dear")
	return true
}

func (self *PlayerLogin) doAction(player *PlayerMgr.Player) bool {
	loger.Debug("My Friend.")
	return true
}

func (self *PlayerLogin) ProcessAction(player *PlayerMgr.Player) bool {
	if false == self.checkAction(player) { //! 检查参数
		return false
	}

	if false == self.payAction(player) { //! 支付代价
		return false
	}

	if false == self.doAction(player) { //! 执行行为
		return false
	}

	return true
}
