package module

import (
	"db"
	"gameserver/table"
	"loger"
)

type PlayerMap map[int64]*Player
type PlayerInfoMap map[int64]*PlayerInfo
type PlayerNameMap map[string]int64

//! 玩家管理器
type PlayerMgr struct {
	moduleMgr     *ModuleMgr    //! 模块管理器指针
	playerMap     PlayerMap     //! 玩家连接表
	playerInfoMap PlayerInfoMap //! 玩家信息表
	playerNameMap PlayerNameMap //! 玩家姓名表
	playerCount   int           //! 当前在线玩家数量
	countLimit    int           //! 人数限制
}

//! 初始化玩家管理器
func (self *PlayerMgr) Init(limit int) {
	self.playerMap = make(PlayerMap)
	self.playerInfoMap = make(PlayerInfoMap)
	self.countLimit = limit

	//! 初始化消息分拣器
	G_Dispatch.Init()
}

//! 设置模块管理器指针
func (self *PlayerMgr) SetModuleMgrPoint(moduleMgr *ModuleMgr) {
	self.moduleMgr = moduleMgr
}

//! 获取模块管理器指针
func (self *PlayerMgr) GetModuleMgrPoint() *ModuleMgr {
	return self.moduleMgr

}

//! 获取当期在线玩家数量
func (self *PlayerMgr) GetOnlinePlayerCount() int {
	return self.playerCount
}

//! 生成新玩家ID
func (self *PlayerMgr) CreateNewPlayerID() int64 {
	playerLst := []PlayerInfo{}
	var playerID int64
	db.Find_Sort(table.GameDB, table.PlayerInfoTable, "_id", -1, 1, &playerLst)
	if len(playerLst) <= 0 {
		playerID = 1
	} else {
		playerID = playerLst[0].AccountID + 1
	}
	return playerID
}

//! 添加一个玩家
func (self *PlayerMgr) AddPlayerCount(player *Player) {

	//! 检测人数是否超标
	if self.countLimit <= self.playerCount {
		loger.Warn("Server limit: %d  current number: %d  can't add new player!", self.countLimit, self.playerCount)
		return
	}

	//! 加入玩家信息表
	self.playerMap[player.PlayerID] = player

	//! 玩家人数 + 1
	self.playerCount += 1
}

//! 减去一个玩家
func (self *PlayerMgr) SubPlayerCount(player *Player) {

	//! 删除该玩家
	delete(self.playerMap, player.PlayerID)

	//! 玩家人数 - 1
	self.playerCount -= 1
}

//! 添加一个玩家信息
func (self *PlayerMgr) AddPlayerInfo(player *PlayerInfo) {
	self.playerInfoMap[player.PlayerID] = player
}

//! 删除一个玩家信息
func (self *PlayerMgr) SubPlayerInfo(player *PlayerInfo) {
	delete(self.playerInfoMap, player.PlayerID)
}

//! 获取一个玩家信息
func (self *PlayerMgr) GetPlayerInfo(playerID int64) *PlayerInfo {
	//! 判断该玩家是否存在
	_, ok := self.playerInfoMap[playerID]
	if ok == false {
		player := new(PlayerInfo)

		//! 玩家不存在,尝试从数据库取出数据
		isExist := db.Find(table.GameDB, table.PlayerInfoTable, "_id", playerID, player)
		if isExist == false {
			//! 玩家信息不存在
			return nil
		}

		//! 取出玩家信息,加入内存
		self.AddPlayerInfo(player)
	}

	return self.playerInfoMap[playerID]
}

//! 获取一个玩家套接字信息
func (self *PlayerMgr) GetPlayerSocket(playerID int64) *Player {
	s, ok := self.playerMap[playerID]
	if ok == false {
		return nil
	}
	return s
}

//! 加入一个玩家信息到数据库
func (self *PlayerMgr) AddPlayerInfoToDB(info *PlayerInfo) bool {
	isSuccess := db.Insert(table.GameDB, table.PlayerInfoTable, info)
	if isSuccess == false {
		return false
	}

	self.playerInfoMap[info.PlayerID] = info
	return true
}

//! 通过账号ID取得一个玩家信息
func (self *PlayerMgr) GetPlayerInfoFromAccount(accountID int64) *PlayerInfo {
	var ret *PlayerInfo

	//! 遍历玩家表,匹配账户ID
	for _, v := range self.playerInfoMap {
		if v.AccountID == accountID {
			ret = v
			break
		}
	}

	if ret == nil {
		//! 玩家不存在,尝试从数据库取出数据
		ret = new(PlayerInfo)
		isExist := db.Find(table.GameDB, table.PlayerInfoTable, "accountid", accountID, ret)

		if isExist == false {
			//! 查无此人
			return nil
		}

		//! 加入玩家表
		self.AddPlayerInfo(ret)
	}
	return ret
}

//! 踢出一个在线玩家 From ID
func (self *PlayerMgr) kickPlayerFromID(playerID int64) {

	//! 根据ID获取玩家信息
	player, ok := self.playerMap[playerID]
	if ok == false {
		loger.Error("Can't find player id: %v", playerID)
		return
	}

	player.ws.Close()

	self.SubPlayerCount(player)
}

func NewPlayerMgr(limit int) *PlayerMgr {
	mgr := new(PlayerMgr)
	mgr.Init(limit)
	return mgr
}
