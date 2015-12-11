package msg

import (
	"gameserver/playermgr"
)

type PlayerLogin_Res struct {
	PlayerMgr.MsgHead    `json:"head"` //! 消息头
	StatusCode           int           `json:"code"` //! 返回码
	PlayerMgr.PlayerInfo `json:"info"` //! 玩家基本信息
}

func (self *PlayerLogin_Res) GetTypeAndAction() (string, string) {
	return "login", "login_result"
}

type CreateRole_Res struct {
	PlayerMgr.MsgHead    `json:"head"` //! 消息头
	StatusCode           int           `json:"code"` //! 返回码
	PlayerMgr.PlayerInfo `json:"info"` //! 玩家基本信息
}

func (self *CreateRole_Res) GetTypeAndAction() (string, string) {
	return "login", "create_result"
}
