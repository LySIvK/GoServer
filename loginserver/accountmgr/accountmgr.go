package AccountMgr

import (
	"db"
	"loger"
	g "loginserver/gameservermgr"
	"loginserver/table"
	"sync"
	"time"
	"tool"
)

//! 账号信息结构体
type AccountInfo struct {
	AccountID      int64  `bson:"_id"`            //! 账户ID
	Name           string `json:"name"`           //! 账户用户名
	Password       string `json:"password"`       //! 账户密码
	CreateTime     int64  `json:"createtime"`     //! 账号创建时间
	LastLoginTime  int64  `json:"lastlogintime"`  //! 上次登录时间
	LoginDays      int    `json:"logindays"`      //! 连续登陆天数
	DeviceID       int    `json:"deviceid"`       //! 设备唯一ID
	Forbidden      bool   `json:"forbidden"`      //! 是否禁止登陆
	LastServerID   int    `json:"lastserverid"`   //! 上次选择的服务器ID
	LoginServerIDs []int  `json:"loginserverids"` //! 登录过的游戏服务器ID
}

//! 账号管理器
type AccountMgr struct {
	gameServerMgr  *g.GameServerMgr       //! 游戏服务器管理器
	lock           sync.RWMutex           //! 协程读写锁
	curAccountID   int64                  //! 当前账号ID
	accountNameMap map[string]int64       //! 账户用户名表
	accountInfoMap map[int64]*AccountInfo //! 账户信息表
	loginKeyMap    map[int64]string       //! 登录键值表
}

//! 初始化
func (self *AccountMgr) Init(gameServerMgr *g.GameServerMgr) {
	self.accountInfoMap = make(map[int64]*AccountInfo)
	self.accountNameMap = make(map[string]int64)
	self.loginKeyMap = make(map[int64]string)
	self.gameServerMgr = gameServerMgr

	//! 获取当前账户ID
	lastAccount := []AccountInfo{}
	db.Find_Sort(table.AccountDB, table.AccountInfoTable, "_id", -1, 1, &lastAccount)
	if len(lastAccount) <= 0 {
		self.curAccountID = 1
	} else {
		self.curAccountID = lastAccount[0].AccountID + 1
	}

	//! 预取一周之内登录过的用户
	prefetching := []AccountInfo{}
	now := time.Now().Unix()
	beginTime := now - (60 * 60 * 24 * 7)
	db.Find_Range(table.AccountDB, table.AccountInfoTable, "lastlogintime", beginTime, now, true, &prefetching)

	if len(prefetching) <= 0 {
		loger.Debug("Prefetching done, but it's zero. time: %v --- %v", beginTime, now)
		return
	}

	for _, v := range prefetching {
		self.accountInfoMap[v.AccountID] = &v
		self.accountNameMap[v.Name] = v.AccountID
	}
}

//! 根据玩家ID获取账户信息
func (self *AccountMgr) GetAccountInfoFromID(id int64) *AccountInfo {
	//! 设置读锁
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.accountInfoMap[id]
}

//! 根据玩家姓名获取账户信息
func (self *AccountMgr) GetAccountInfoFromName(name string) *AccountInfo {
	//! 设置读锁
	self.lock.RLock()
	defer self.lock.RUnlock()

	playerID := self.accountNameMap[name]
	if playerID == 0 {
		return nil
	}

	return self.accountInfoMap[playerID]
}

//! 检测名字是否重复
func (self *AccountMgr) IsNameExist(name string) bool {
	//! 设置读锁
	self.lock.RLock()
	defer self.lock.RUnlock()

	playerID := self.accountNameMap[name]
	if playerID == 0 {
		//! 内存中没有,去数据库中查找
		info := AccountInfo{}
		err := db.Find(table.AccountDB, table.AccountInfoTable, "name", name, &info)
		if err != nil {
			return false
		}

		//! 找到后加入缓存
		self.lock.Lock()
		self.accountNameMap[info.Name] = info.AccountID
		self.accountInfoMap[info.AccountID] = &info
		self.lock.Unlock()
	}

	return true
}

//! 获取注册帐号ID
func (self *AccountMgr) CreateNewAccountID() int64 {
	//! 多开登陆服时,每次读取数据库取出最新AccountID
	lastAccount := []AccountInfo{}
	db.Find_Sort(table.AccountDB, table.AccountInfoTable, "_id", -1, 1, &lastAccount)
	if len(lastAccount) <= 0 {
		self.curAccountID = 1
	} else {
		self.curAccountID = lastAccount[0].AccountID + 1
	}

	ret := self.curAccountID
	self.curAccountID += 1
	return ret
}

//! 生成新账号
func (self *AccountMgr) CreateNewAccountInfo(accounName string, accountPwd string, deviceID int) *AccountInfo {
	info := new(AccountInfo)
	info.AccountID = self.CreateNewAccountID()
	info.Name = accounName
	info.Password = tool.MD5(accountPwd)
	info.CreateTime = time.Now().Unix()
	info.LastLoginTime = time.Now().Unix()
	info.LoginDays = 1
	info.DeviceID = deviceID
	info.Forbidden = false
	info.LastServerID = 0
	info.LoginServerIDs = []int{}
	return info
}

//! 添加登陆Key
func (self *AccountMgr) AddLoginKey(accountID int64, key string) {
	//! 设置写锁
	self.lock.Lock()
	defer self.lock.Unlock()

	self.loginKeyMap[accountID] = key
}

//! 检测登陆Key
func (self *AccountMgr) CheckLoginKey(accountID int64, key string) bool {
	//! 设置读锁
	self.lock.RLock()
	defer self.lock.RUnlock()

	loginKey, ok := self.loginKeyMap[accountID]
	if ok == true && key == loginKey {
		return true
	}

	return false
}

//! 加入账号列表
func (self *AccountMgr) AddAccountInfo(info *AccountInfo) {
	//! 设置写锁
	self.lock.Lock()
	defer self.lock.Unlock()

	self.accountNameMap[info.Name] = info.AccountID
	self.accountInfoMap[info.AccountID] = info
}

func NewAccountMgr(gameServerMgr *g.GameServerMgr) *AccountMgr {
	mgr := new(AccountMgr)
	mgr.Init(gameServerMgr)
	return mgr
}
