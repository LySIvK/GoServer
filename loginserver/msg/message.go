package msg

//! 游戏请求注册消息请求
type Msg_RegisterGameServer_Req struct {
	ServerID   int    `json:"serverid"`   //! 游戏服务器ID
	ServerName string `json:"servername"` //! 游戏服务器名字
	ServerIP   string `json:"serverip"`   //! 游戏服务器IP
	PlayerNum  int    `json:"playernum"`  //! 当前人数
	IsNew      bool   `json:"new"`        //! 是否为新服
}

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

//! 用户注册请求
type Msg_UserRegister_Req struct {
	AccountName string `json:table.AccountInfoTable` //! 注册帐号
	AccountPwd  string `json:"password"`             //! 注册密码
}

//! 用户注册请求返回
type Msg_UserRegister_Res struct {
	StatusCode int `json:"statuscode"` //! 状态码
}

//! 服务器列表请求
type Msg_ServerList_Req struct {
	AccountID int64 `json:"accountid"` //! 账户ID
}

//! 服务器列表请求返回
type GameServerInfo struct {
	ID         int    `json:"id"`         //! 游戏服务器ID
	Name       string `json:"name"`       //! 游戏服务器名字
	PlayerNum  int    `json:"playernum"`  //! 当前游戏人数
	Status     int    `json:"status"`     //! 游戏服务器状态 0-> 关闭  1-> 正常 2-> 无心跳...
	UpdateTime int64  `json:"updatetime"` //! 更新时间
	Addr       string `json:"addr"`       //! 游戏服务器地址+端口
	IsNew      bool   `json:"new"`        //! 是否为新服
}

type Msg_ServerList_Res struct {
	StatusCode int              `json:"statuscode"` //! 状态码
	ServerLst  []GameServerInfo `json:"serverlst"`  //! 服务器信息列表
}

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
