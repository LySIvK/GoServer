package PlayerMgr

import ()

type PlayerInfo struct {
	PlayerID  int64   `bson:"_id"`       //! 玩家唯一标识
	AccountID int64   `json:"accountid"` //! 玩家账号ID
	Name      string  `json:"name"`      //! 昵称
	Level     int     `json:"level"`     //! 等级
	Skill     []int   `json:"skill"`     //! 技能
	Money     []int64 `json:"money"`     //! 银币
	GuildID   int     `json:"guildid"`   //! 公会ID
}
