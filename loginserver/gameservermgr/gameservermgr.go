package GameServerMgr

import (
	"loger"
	"sync"
	"time"
)

//! 游戏服务器信息结构体
type GameServerInfo struct {
	ID         int    `json:"id"`         //! 游戏服务器ID
	Name       string `json:"name"`       //! 游戏服务器名字
	PlayerNum  int    `json:"playernum"`  //! 当前游戏人数
	Status     int    `json:"status"`     //! 游戏服务器状态 0-> 关闭  1-> 正常 2-> 无心跳...
	UpdateTime int64  `json:"updatetime"` //! 更新时间
	Addr       string `json:"addr"`       //! 游戏服务器地址+端口
	IsNew      bool   `json:"new"`        //! 是否为新服
}

type GameServerInfoMap map[int]*GameServerInfo

//! 游戏服务器管理器
type GameServerMgr struct { //! WARNNING: 私有变量,外部不可访问
	serverMap GameServerInfoMap
	lock      sync.RWMutex
}

//! 得到当前服务器个数
func (self *GameServerMgr) GetServerCount() int {
	return len(self.serverMap)
}

//! 初始化游戏服务器管理器
func (self *GameServerMgr) Init() {
	//! 存储服务器map
	self.serverMap = make(GameServerInfoMap)

	//! 开启循检
	go self.Inspect()
}

//! 循检服务器状态
func (self *GameServerMgr) Inspect() {

	//! 30秒循检一次状态
	inspectTime := time.Tick(30 * time.Second)
	for {
		curTime := time.Now().Unix()
		for _, info := range self.serverMap {
			if curTime-info.UpdateTime > 30 {
				//! 30秒未有心跳包则判定服务器断线
				self.lock.Lock()
				info.Status = 2
				loger.Warn("GameServer was NO response! ServerID: %d ServerName: %v ServerAddr: %v",
					info.ID, info.Name, info.Addr)
				self.lock.Unlock()
			}
		}

		<-inspectTime
	}
}

//! 增加游戏服务器信息
func (self *GameServerMgr) AddNewGameServerInfo(info *GameServerInfo) {
	self.lock.Lock()
	defer self.lock.Unlock()

	pCurServer := self.serverMap[info.ID]
	if pCurServer != nil {
		//! 已存在该服务器信息
		info.UpdateTime = time.Now().Unix()
		pCurServer.Status = 1
		pCurServer.UpdateTime = info.UpdateTime
	} else {
		//! 不存在该服务器信息
		self.serverMap[info.ID] = info
	}
}

//! 获取游戏服务器名称
func (self *GameServerMgr) GetGameServerName(serverID int) string {
	if self.serverMap[serverID] != nil {
		return self.serverMap[serverID].Name
	}
	return ""
}

//! 获取游戏服务器地址
func (self *GameServerMgr) GetGameServerAddr(serverID int) string {
	if self.serverMap[serverID] != nil {
		return self.serverMap[serverID].Addr
	}
	return ""
}

//! 获取游戏服务器属性
func (self *GameServerMgr) GetGameServerInfo(serverID int) *GameServerInfo {
	return self.serverMap[serverID]
}

//! 获取游戏服务器个数
func (self *GameServerMgr) GetGameServerCount() int {
	return len(self.serverMap)
}

//! 获取连接状态服务器个数
func (self *GameServerMgr) GetGameServerConnectCount() int {
	serverCount := 0
	for _, v := range self.serverMap {
		if v.Status == 1 {
			serverCount++
		}
	}
	return serverCount
}

//! 获取异常状态服务器列表
func (self *GameServerMgr) GetGameServerWarnLst() []*GameServerInfo {
	var lst []*GameServerInfo
	for _, v := range self.serverMap {
		if v.Status != 1 {
			lst = append(lst, v)
		}
	}
	return lst
}

//! 获取所有服务器列表
func (self *GameServerMgr) GetGameServerLst() []*GameServerInfo {
	var lst []*GameServerInfo
	for _, v := range self.serverMap {
		lst = append(lst, v)
	}
	return lst
}

//! 获取推荐服务器
func (self *GameServerMgr) GetRecommendServerID() *GameServerInfo {
	self.lock.RLock()
	defer self.lock.RUnlock()

	serverID := 0
	for _, v := range self.serverMap {
		if v.ID > serverID && v.Status == 1 {
			serverID = v.ID
		}
	}

	if serverID == 0 {
		return nil
	}

	return self.serverMap[serverID]
}

//! 生成游戏服务器管理器
func NewGameServerMgr() *GameServerMgr {
	mgr := new(GameServerMgr)
	mgr.Init()
	return mgr
}
