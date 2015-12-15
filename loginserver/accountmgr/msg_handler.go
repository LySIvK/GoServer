package AccountMgr

import (
	"db"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"loger"
	"loginserver/msg"
	"loginserver/table"
	"net/http"
	"time"
	"tool"
)

//! 消息处理-玩家登陆
func (self *AccountMgr) Handler_UserLogin(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 收到玩家登陆消息
	loger.Debug("Recv msg from %v", r.URL.String())

	var req msg.Msg_UserLogin_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		loger.Error("Handler_UserLogin Unmarshal error: %s", err.Error())
		return
	}

	//! 创建响应消息
	var response msg.Msg_UserLogin_Res
	response.StatusCode = msg.RE_UNKNOW_ERR

	//! 设置延后发送
	defer func() {
		b, err := json.Marshal(&response)
		if err != nil {
			loger.Error("Handler_UserLogin Marshal error: %s", err.Error())
			return
		}

		loger.Debug("Send: %s", string(b))
		w.Write(b)
	}()

	//! 参数检查
	result, errCode := self.Handler_UserLogin_Check(req.AccountName, req.AccountPwd)
	if result != true {
		response.StatusCode = errCode
		return
	}

	//! 读取用户信息
	info := self.GetAccountInfoFromName(req.AccountName)
	response.StatusCode = msg.RE_SUCCESS
	response.AccountID = info.AccountID
	response.LoginKey = bson.NewObjectId().Hex()
	response.LastLoginServerID = info.LastServerID
	if info.LastServerID == 0 {
		//! 若无上次登陆服,则获取推荐服
		recommendServer := self.gameServerMgr.GetRecommendServerID()
		if recommendServer != nil {
			response.LastLoginServerID = recommendServer.ID
			response.LastLoginServerName = recommendServer.Name
			response.LastLoginServerAddr = recommendServer.Addr
		}
	}

	if info.LastServerID > 0 {
		response.LastLoginServerName = self.gameServerMgr.GetGameServerName(info.LastServerID)
		response.LastLoginServerAddr = self.gameServerMgr.GetGameServerAddr(info.LastServerID)
	}
	//! 将登陆Key加入缓存
	self.AddLoginKey(info.AccountID, response.LoginKey)

	//! 检测连续登陆
	now := time.Now().Unix()
	if now-info.LastLoginTime <= 24*60*60 {
		info.LoginDays += 1
		go db.IncFieldValue(table.AccountDB, table.AccountInfoTable, "_id", info.AccountID, "logindays", 1)
	} else {
		info.LoginDays = 1
		go db.UpdateField(table.AccountDB, table.AccountInfoTable, "_id", info.AccountID, "logindays", 1)
	}
}

//! 玩家登陆消息检测
func (self *AccountMgr) Handler_UserLogin_Check(accountName string, accountPwd string) (bool, int) {
	//! 检查用户名长度
	if len(accountName) < 6 || len(accountName) > 32 {
		loger.Warn("Account name is invalid. AccountName: %s", accountName)

		return false, msg.RE_INVALID_ACCOUNTNAME
	}

	//! 检查密码长度
	if len(accountPwd) < 6 || len(accountPwd) > 32 {
		loger.Warn("Account password is invalid. Password: %s", accountPwd)

		return false, msg.RE_INVALID_PASSWORD
	}

	//! 检测用户名是否存在
	bCheck := self.IsNameExist(accountName)
	if bCheck != true {
		loger.Error("Account name is not exists. AccountName: %s", accountName)

		return false, msg.RE_ACCOUNT_NOT_EXIST
	}

	//! 检测用户密码是否正确
	info := self.GetAccountInfoFromName(accountName)
	if info.Password != tool.MD5(accountPwd) {
		return false, msg.RE_INVALID_PASSWORD
	}

	return true, msg.RE_SUCCESS
}

//! 用户注册
func (self *AccountMgr) Handler_UserRegister(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 收到玩家信息
	loger.Debug("Recv msg from %v", r.URL.String())

	//! 解析消息
	var req msg.Msg_UserRegister_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		loger.Error("Handler_UserRegister Unmarshal error: %s", err.Error())
		return
	}

	//! 创建回应消息
	var response msg.Msg_UserRegister_Res
	response.StatusCode = msg.RE_UNKNOW_ERR

	//! 设置延后发送
	defer func() {
		res, err := json.Marshal(&response)
		if err != nil {
			loger.Error("Handler_UserLogin Marshal error: %s", err.Error())
			return
		}

		loger.Debug("Send: %s", string(res))
		w.Write(res)
	}()

	//! 检测参数
	ok, errCode := self.Handler_UserRegister_Check(req.AccountName, req.AccountPwd)
	if ok == false {
		response.StatusCode = errCode
		return
	}

	//! 注册帐号
	newInfo := self.CreateNewAccountInfo(req.AccountName, req.AccountPwd, 0)
	self.AddAccountInfo(newInfo)

	isSuccess := db.Insert(table.AccountDB, table.AccountInfoTable, newInfo)
	if isSuccess == true {
		response.StatusCode = msg.RE_SUCCESS
	}
}

//! 注册检测
func (self *AccountMgr) Handler_UserRegister_Check(accountName string, accountPwd string) (bool, int) {
	//! 检查用户名长度
	if len(accountName) < 6 || len(accountName) > 32 {
		loger.Warn("Account name is invalid. AccountName: %s", accountName)

		return false, msg.RE_INVALID_ACCOUNTNAME
	}

	//! 检查密码长度
	if len(accountPwd) < 6 || len(accountPwd) > 32 {
		loger.Warn("Account password is invalid. Password: %s", accountPwd)

		return false, msg.RE_INVALID_PASSWORD
	}

	//! 检测用户名是否存在
	bCheck := self.IsNameExist(accountName)
	if bCheck == true {
		loger.Error("Account name already exists. AccountName: %s", accountName)
		return false, msg.RE_ACCOUNT_EXIST
	}

	return true, msg.RE_SUCCESS
}

//! 请求服务器列表
func (self *AccountMgr) Handler_ServerList(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 输出用户消息
	loger.Debug("recv msg from %v", r.URL.String())

	var req msg.Msg_ServerList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		loger.Error("Handler_ServerList Unmarshal error: %s", err.Error())
		return
	}

	//! 创建回应消息
	var res msg.Msg_ServerList_Res
	res.StatusCode = msg.RE_UNKNOW_ERR

	defer func() {
		b, err := json.Marshal(&res)
		if err != nil {
			loger.Error("Handler_ServerList Marshal error: %s", err.Error())
			return
		}

		loger.Debug("Send: %s", string(b))
		w.Write(b)
	}()

	//! 检查参数
	ok, errCode := self.Handler_ServerList_Check(req.AccountID)
	if ok == false {
		res.StatusCode = errCode
		return
	}

	serverLst := self.gameServerMgr.GetGameServerLst()
	for _, v := range serverLst {
		var data msg.GameServerInfo
		data.ID = v.ID
		data.Addr = v.Addr
		data.IsNew = v.IsNew
		data.Name = v.Name
		data.PlayerNum = v.PlayerNum
		data.Status = v.Status
		data.UpdateTime = v.UpdateTime
		res.ServerLst = append(res.ServerLst, data)
	}

	res.StatusCode = msg.RE_SUCCESS
}

//! 请求服务器列表参数检查
func (self *AccountMgr) Handler_ServerList_Check(accountID int64) (bool, int) {
	_, ok := self.loginKeyMap[accountID]
	if ok == false {
		return false, msg.RE_NOT_LOGIN
	}

	return true, msg.RE_SUCCESS
}

//! 验证用户登陆
func (self *AccountMgr) Handler_VerifyUserLogin(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	loger.Debug("recv msg from %v", r.URL.String())

	//! 解析消息
	var req msg.Msg_VerifyUserLogin_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		loger.Error("Handler_VerifyUserLogin Unmarshal error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.Msg_VerifyUserLogin_Res
	response.StatusCode = msg.RE_UNKNOW_ERR

	defer func() {
		b, err := json.Marshal(&response)
		if err != nil {
			loger.Error("Handler_VerifyUserLogin Marshal error: %s", err.Error())
			return
		}
		loger.Debug("Send: %v", string(b))
		w.Write(b)
	}()

	//! 检查参数
	ok, errCode := self.Handler_VerifyUserLogin_Check(req.AccountID, req.LoginKey)
	if ok == false {
		response.StatusCode = errCode
		return
	}

	//! 获取用户信息
	accountInfo := self.GetAccountInfoFromID(req.AccountID)
	if accountInfo == nil {
		loger.Error("GetAccountInfo Fail. AccountID: %v", req.AccountID)
		return
	}

	accountInfo.LastServerID = req.ServerID
	isExist := false
	for _, v := range accountInfo.LoginServerIDs {
		if v == req.ServerID {
			isExist = true
			break
		}
	}

	if isExist == false {
		//! 如果不存在,则加入用户登陆过的服务器列表
		accountInfo.LoginServerIDs = append(accountInfo.LoginServerIDs, req.ServerID)
		go db.AddToArray(table.AccountDB, table.AccountInfoTable, "_id", req.AccountID, "loginserverids", req.ServerID)
	}

	//! 设置用户最后登录服务器
	go db.UpdateField(table.AccountDB, table.AccountInfoTable, "_id", req.AccountID, "lastserverid", req.ServerID)

	response.StatusCode = msg.RE_SUCCESS
}

//! 验证用户登陆检查参数
func (self *AccountMgr) Handler_VerifyUserLogin_Check(accountID int64, key string) (bool, int) {
	result := self.CheckLoginKey(accountID, key)
	if result == false {
		return false, msg.RE_NOT_LOGIN
	}

	return true, msg.RE_SUCCESS
}
