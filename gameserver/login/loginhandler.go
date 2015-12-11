package login

import (
	"bytes"
	"db"
	"encoding/json"
	"fmt"
	retmsg "gameserver/msg"
	"gameserver/playermgr"
	"gameserver/table"
	"loger"
	"loginserver/msg"
	"net/http"
	"serverconfig"
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
	if self.AccountID <= 0 {
		loger.Error("AccountID is invalid. AccountID: %v", self.AccountID)
		return false
	}

	if len(self.LoginKey) <= 0 {
		loger.Error("LoginKey is invalid. LoginKey: %s", self.LoginKey)
		return false
	}

	return true
}

func (self *PlayerLogin) payAction(player *PlayerMgr.Player, res *msg.Msg_VerifyUserLogin_Res) bool {
	//! 请求登录服务器,验证该用户是否登录
	loginServerUrl := fmt.Sprintf("http://%s:%d", serverconfig.G_Config.LoginServer_IP, serverconfig.G_Config.LoginServer_Port)
	loginServerUrl += "/verifyuserlogin"

	//! 构建请求
	var req msg.Msg_VerifyUserLogin_Req
	req.AccountID = self.AccountID
	req.LoginKey = self.LoginKey
	req.ServerID = serverconfig.G_Config.GameServerID

	//! 解析为Json格式
	b, err := json.Marshal(&req)
	if err != nil {
		loger.Error("Marshal fail. Error: %s", err.Error())
		return false
	}

	//! 发送请求
	resp, err := http.Post(loginServerUrl, "Text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("Http post fail. Error: %s", err.Error())
		return false
	}

	//! 获取结果
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()

	//! 解析到结构体
	err = json.Unmarshal(buffer, res)
	if err != nil {
		loger.Error("Unmarshal fail. Error: %s", err.Error())
		return false
	}

	return true
}

func (self *PlayerLogin) doAction(player *PlayerMgr.Player, res *msg.Msg_VerifyUserLogin_Res) bool {
	//! 创建返回信息
	msg := new(retmsg.PlayerLogin_Res)
	msg.StatusCode = retmsg.RE_UNKNOW_ERR

	defer player.Send(msg)

	//! 判断返回码
	if res.StatusCode != retmsg.RE_SUCCESS {
		//! 玩家没有登录
		loger.Error("Verify fail. StatusCode: %d", res.StatusCode)
		msg.StatusCode = retmsg.RE_PLAYER_NOT_LOGIN
		return false
	}

	//! 判断数据库中是否存在该玩家信息
	isExist := db.IsRecordExist(table.GameDB, table.PlayerInfoTable, "accountid", self.AccountID)
	if isExist == true {
		//! 获取玩家管理器指针
		playerMgr := player.PlayerMgr

		//! 读取玩家信息
		info := playerMgr.GetPlayerInfoFromAccount(self.AccountID)
		if info != nil {
			msg.StatusCode = retmsg.RE_SUCCESS
			msg.PlayerInfo = *info
			return true
		}
	}

	msg.StatusCode = retmsg.RE_ROLE_NOT_EXIST
	return true
}

func (self *PlayerLogin) ProcessAction(player *PlayerMgr.Player) bool {
	if false == self.checkAction(player) { //! 检查参数
		return false
	}

	var res msg.Msg_VerifyUserLogin_Res
	if false == self.payAction(player, &res) { //! 支付代价
		return false
	}

	if false == self.doAction(player, &res) { //! 执行行为
		return false
	}

	return true
}
