package main

import (
	"bytes"
	"encoding/json"
	"loger"
	"net/http"
)

var (
	regurl       = "http://192.168.1.106:9016/register"
	loginurl     = "http://192.168.1.106:9016/login"
	verifyurl    = "http://192.168.1.106:9016/verifyuserlogin"
	serverlsturl = "http://192.168.1.106:9016/serverlist"
)

func TestRegister() {
	//! 用户注册请求
	type Msg_UserRegister_Req struct {
		AccountName string `json:table.AccountInfoTable` //! 注册帐号
		AccountPwd  string `json:"password"`             //! 注册密码
	}

	//! 用户注册请求返回
	type Msg_UserRegister_Res struct {
		StatusCode int `json:"statuscode"` //! 状态码
	}

	var req Msg_UserRegister_Req
	req.AccountName = "test001"
	req.AccountPwd = "000000"

	b, _ := json.Marshal(&req)
	resp, err := http.Post(regurl, "text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("post error: %s", err.Error())
		return
	}

	var res Msg_UserRegister_Res
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()
	json.Unmarshal(buffer, &res)

	loger.Info("RetCode: %d", res.StatusCode)
}

func TestLogin() {
	//! 用户登陆请求
	type Msg_UserLogin_Req struct {
		AccountName string `json:table.AccountInfoTable` //! 玩家帐号
		AccountPwd  string `json:"password"`             //! 玩家密码
	}

	//! 用户登陆请求返回
	type Msg_UserLogin_Res struct {
		StatusCode          int    `json:"statuscode"`     //! 状态码
		AccountID           int64  `json:"accountid"`      //! 帐号ID
		LoginKey            string `json:"key"`            //! 登陆键值
		LastLoginServerID   int    `json:"lastserverid"`   //! 上次登陆服务器ID
		LastLoginServerName string `json:"lastservername"` //! 上次登陆服务器名字
		LastLoginServerAddr string `json:"lastserveraddr"` //! 上次登陆服务器地址
	}

	var req Msg_UserLogin_Req
	req.AccountName = "test001"
	req.AccountPwd = "000000"

	b, _ := json.Marshal(&req)
	resp, err := http.Post(loginurl, "Text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("post error: %s", err.Error())
		return
	}

	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()

	loger.Info("recv: %s", string(buffer))
}

func TestVerifyLogin() {
	//! 验证用户登陆请求
	type Msg_VerifyUserLogin_Req struct {
		AccountID int64  `json:"id"`       //! 账号ID
		LoginKey  string `json:"key"`      //! 登陆键值
		ServerID  int    `json:"serverid"` //! 服务器ID
	}

	//! 验证用户登陆请求返回
	type Msg_VerifyUserLogin_Res struct {
		StatusCode int `json:"statuscode"` //! 状态码
	}

	var req Msg_VerifyUserLogin_Req
	req.AccountID = 1
	req.LoginKey = "566930d5aeddbf2b28d277aa"
	req.ServerID = 1

	b, _ := json.Marshal(&req)
	resp, err := http.Post(verifyurl, "Text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("post error: %s", err.Error())
		return
	}

	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()

	loger.Info("recv: %s", string(buffer))
}

func TestServerLst() {
	type Msg_ServerList_Req struct {
		AccountID int64 `json:"accountid"` //! 账户ID
	}

	var req Msg_ServerList_Req
	req.AccountID = 1
	b, _ := json.Marshal(&req)
	resp, err := http.Post(serverlsturl, "Text/HTML", bytes.NewReader(b))
	if err != nil {
		loger.Error("post error: %s", err.Error())
		return
	}

	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()

	loger.Info("recv: %s", string(buffer))
}
