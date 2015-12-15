package module

import (
	"gameserver/msg"
	"loger"
)

//! 聊天消息请求
type Msg_Chat_Private struct {
	MsgHead    `json:"head"`
	SendUserID int64  `json:"sendid"` //! 私聊ID
	Message    string `json:"msg"`    //! 发送内容
}

func (self *Msg_Chat_Private) GetTypeAndAction() (string, string) {
	return "chat", "private"
}

func (self *Msg_Chat_Private) checkAction(player *Player, res *Chat_Private_Res) bool {
	//! 检测玩家登陆
	if player.PlayerID <= 0 {
		loger.Error("Player not login.")
		res.StatusCode = msg.RE_NOT_LOGIN
		return false
	}

	//! 检测消息长度
	if len(self.Message) > 512 {
		loger.Error("Chat message too long. PlayerID: %v", player.PlayerID)
		res.StatusCode = msg.RE_TOO_LONG_CHAT
		return false
	}

	//！检测私聊用户ID
	if self.SendUserID <= 0 {
		loger.Error("Chat message send user id is zero. PlayerID: %v", player.PlayerID)
		res.StatusCode = msg.RE_INVALID_PLAYERID
		return false
	}

	//! 检测私聊用户在线
	info := player.PlayerMgr.GetPlayerInfo(self.SendUserID)
	if info == nil {
		loger.Error("Send user is not online. PlayerID: %v", player.PlayerID)
		res.StatusCode = msg.RE_SEND_PLAYER_NOT_ONLINE
		return false
	}

	return true
}

func (self *Msg_Chat_Private) payAction(player *Player, res *Chat_Private_Res) bool {
	//! TODO: 支付处理
	return true
}

func (self *Msg_Chat_Private) doAction(player *Player, res *Chat_Private_Res) bool {
	//! 获取玩家套接字
	playerMgr := player.PlayerMgr
	socket := playerMgr.GetPlayerSocket(player.PlayerID)

	//! 创建封包
	newMsg := new(Chat_Private_Send)
	newMsg.Message = self.Message

	//! 发送封包
	socket.Send(newMsg)

	return true
}

func (self *Msg_Chat_Private) ProcessAction(player *Player) bool {

	res := new(Chat_Private_Res)
	res.StatusCode = msg.RE_UNKNOW_ERR
	defer player.Send(res)

	if self.checkAction(player, res) == false {
		return false
	}

	if self.payAction(player, res) == false {
		return false
	}

	if self.doAction(player, res) == false {
		return false
	}

	res.StatusCode = msg.RE_SUCCESS
	return true
}
