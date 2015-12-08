package GameServerMgr

import (
	"encoding/json"
	"loger"
	"loginserver/msg"
	"net/http"
	"time"
)

//! 消息处理-注册游戏服务器
func (self *GameServerMgr) Handler_RegisterGameSvr(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 收到游戏服务器注册信息
	loger.Debug("Recv msg from %v", r.URL.String())

	msg := msg.Msg_RegisterGameServer_Req{}
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		loger.Error("Handler_RegisterGameSvr Unmarshal error: %s", err.Error())
		return
	}

	//! 创建游戏服务器信息
	info := &GameServerInfo{
		ID:         msg.ServerID,
		Name:       msg.ServerName,
		PlayerNum:  msg.PlayerNum,
		Status:     1,
		UpdateTime: time.Now().Unix(),
		Addr:       msg.ServerIP,
		IsNew:      msg.IsNew,
	}

	//! 加入游戏服务器管理器
	self.AddNewGameServerInfo(info)

	loger.Debug("message: %v", string(buffer))
}
