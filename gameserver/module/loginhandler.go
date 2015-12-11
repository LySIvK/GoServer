package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	code "gameserver/msg"
	"loger"
	loginmsg "loginserver/msg"
	"net/http"
	"serverconfig"
)

//! 玩家登录消息
type Msg_PlayerLogin struct {
	MsgHead   `json:"head"`
	AccountID int64
	LoginKey  string
}

func (self *Msg_PlayerLogin) GetTypeAndAction() (string, string) {
	return "login", "login"
}

func (self *Msg_PlayerLogin) checkAction(player *Player, msg *PlayerLogin_Res) bool {
	if self.AccountID <= 0 {
		loger.Error("AccountID is invalid. AccountID: %v", self.AccountID)
		msg.StatusCode = code.RE_INVALID_ACCOUNTID
		return false
	}

	if len(self.LoginKey) <= 0 {
		loger.Error("LoginKey is invalid. LoginKey: %s", self.LoginKey)
		msg.StatusCode = code.RE_INVALID_PASSWORD
		return false
	}

	//! 请求登录服务器,验证该用户是否登录
	loginServerUrl := fmt.Sprintf("http://%s:%d", serverconfig.G_Config.LoginServer_IP, serverconfig.G_Config.LoginServer_Port)
	loginServerUrl += "/verifyuserlogin"

	//! 构建请求
	var req loginmsg.Msg_VerifyUserLogin_Req
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
	var res loginmsg.Msg_VerifyUserLogin_Res

	err = json.Unmarshal(buffer, &res)
	if err != nil {
		loger.Error("Unmarshal fail. Error: %s", err.Error())
		return false
	}

	//! 判断返回码
	if res.StatusCode != code.RE_SUCCESS {
		//! 玩家没有登录
		loger.Error("Verify fail. StatusCode: %d", res.StatusCode)
		msg.StatusCode = code.RE_NOT_LOGIN
		return false
	}

	return true
}

func (self *Msg_PlayerLogin) payAction(player *Player, msg *PlayerLogin_Res) bool {

	return true
}

func (self *Msg_PlayerLogin) doAction(player *Player, msg *PlayerLogin_Res) bool {

	//! 获取玩家管理器指针
	playerMgr := player.PlayerMgr

	//! 读取玩家信息
	info := playerMgr.GetPlayerInfoFromAccount(self.AccountID)
	if info == nil {
		msg.StatusCode = code.RE_ROLE_NOT_EXIST
		return true
	}

	msg.StatusCode = code.RE_SUCCESS
	msg.PlayerInfo = *info
	return true
}

func (self *Msg_PlayerLogin) ProcessAction(player *Player) bool {
	//! 创建返回信息
	msg := new(PlayerLogin_Res)
	msg.StatusCode = code.RE_UNKNOW_ERR
	defer player.Send(msg)

	if false == self.checkAction(player, msg) { //! 检查参数
		return false
	}

	if false == self.payAction(player, msg) { //! 支付代价
		return false
	}

	if false == self.doAction(player, msg) { //! 执行行为
		return false
	}

	return true
}

//! 创建角色信息
type Msg_CreateRole struct {
	MsgHead    `json:"head"`
	AccountID  int64  `json:"accountid"`
	LoginKey   string `json:"key"`
	PlayerName string `json:"name"`
}

func (self *Msg_CreateRole) GetTypeAndAction() (string, string) {
	return "login", "create"
}

func (self *Msg_CreateRole) checkAction(player *Player, msg *CreateRole_Res) bool {
	if self.AccountID <= 0 {
		loger.Error("AccountID is invalid. ID: %v", self.AccountID)
		msg.StatusCode = code.RE_INVALID_ACCOUNTID
		return false
	}

	if len(self.LoginKey) <= 0 {
		loger.Error("LoginKey is invalid. LoginKey: %s", self.LoginKey)
		msg.StatusCode = code.RE_INVALID_PASSWORD
		return false
	}

	if len(self.PlayerName) <= 0 {
		loger.Error("PlayerName is invalid. Name: %v", self.PlayerName)
		msg.StatusCode = code.RE_INVALID_PLAYERNAME
		return false
	}

	//! 请求登录服务器,验证该用户是否登录
	loginServerUrl := fmt.Sprintf("http://%s:%d", serverconfig.G_Config.LoginServer_IP, serverconfig.G_Config.LoginServer_Port)
	loginServerUrl += "/verifyuserlogin"

	//! 构建请求
	var req loginmsg.Msg_VerifyUserLogin_Req
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
	var res loginmsg.Msg_VerifyUserLogin_Res

	err = json.Unmarshal(buffer, &res)
	if err != nil {
		loger.Error("Unmarshal fail. Error: %s", err.Error())
		return false
	}

	//! 判断返回码
	if res.StatusCode != code.RE_SUCCESS {
		//! 玩家没有登录
		loger.Error("Verify fail. StatusCode: %d", res.StatusCode)
		msg.StatusCode = code.RE_NOT_LOGIN
		return false
	}

	//! 判断玩家是否已创建角色
	playerMgr := player.PlayerMgr
	info := playerMgr.GetPlayerInfoFromAccount(self.AccountID)
	if info != nil {
		//! 角色已存在
		loger.Error("Players have created information. playerID: %v  accountID: %v", info.PlayerID, info.AccountID)
		msg.StatusCode = code.RE_ROLE_NOT_EXIST
		return false
	}

	return true
}

func (self *Msg_CreateRole) payAction(player *Player, msg *CreateRole_Res) bool {
	return true
}

func (self *Msg_CreateRole) doAction(player *Player, msg *CreateRole_Res) bool {
	//! 创建角色
	info := new(PlayerInfo)
	info.AccountID = self.AccountID
	info.PlayerID = player.PlayerID
	info.Money = append(info.Money, 5000) //! 初始给予5000金币
	info.Level = 1

	//! 数据库添加一条记录
	playerMgr := player.PlayerMgr
	playerMgr.AddPlayerInfoToDB(info)
	return true
}

func (self *Msg_CreateRole) ProcessAction(player *Player) bool {
	msg := new(CreateRole_Res)
	msg.StatusCode = code.RE_UNKNOW_ERR
	defer player.Send(msg)

	if self.checkAction(player, msg) == false {
		return false
	}

	if self.payAction(player, msg) == false {
		return false
	}

	if self.doAction(player, msg) == false {
		return false
	}
	return true
}
