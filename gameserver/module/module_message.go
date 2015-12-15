package module

import ()

//! 玩家登录消息回应
type PlayerLogin_Res struct {
	MsgHead    `json:"head"` //! 消息头
	StatusCode int           `json:"code"` //! 返回码
	PlayerInfo `json:"info"` //! 玩家基本信息
}

func (self *PlayerLogin_Res) GetTypeAndAction() (string, string) {
	return "login", "login_result"
}

//! 创建角色消息回应
type CreateRole_Res struct {
	MsgHead    `json:"head"` //! 消息头
	StatusCode int           `json:"code"` //! 返回码
	PlayerInfo `json:"info"` //! 玩家基本信息
}

func (self *CreateRole_Res) GetTypeAndAction() (string, string) {
	return "login", "create_result"
}

//! 私聊消息回应
type Chat_Private_Res struct {
	MsgHead    `json:"head"` //! 消息头
	StatusCode int           `json:"code"` //! 返回码
}

func (self *Chat_Private_Res) GetTypeAndAction() (string, string) {
	return "chat", "private"
}

//! 私聊消息发送
type Chat_Private_Send struct {
	MsgHead `json:"head"` //! 消息头
	Message string        `json:"msg"` //! 消息
}

func (self *Chat_Private_Send) GetTypeAndAction() (string, string) {
	return "chat", "private_send"
}
