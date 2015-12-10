package PlayerMgr

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
func (self *PlayerMgr) AddPlayer(player *Player) {

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
func (self *PlayerMgr) SubPlayer(player *Player) {

	//! 删除该玩家
	delete(self.playerMap, player.PlayerID)

	//! 玩家人数 - 1
	self.playerCount -= 1
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

	self.SubPlayer(player)
}

func NewPlayerMgr(limit int) *PlayerMgr {
	mgr := new(PlayerMgr)
	mgr.Init(limit)
	return mgr
}
