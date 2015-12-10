package main

import (
	"db"
	"fmt"
)

var (
	Url       = "127.0.0.1:27017"
	PoolLimit = 50
)

type Player struct {
	ID    int `bson:"_id"`
	Name  string
	Level int
	Money int64
	Skill []int
}

//! 账号信息结构体
type AccountInfo struct {
	AccountID      int64  `bson:"_id"`            //! 账户ID
	Name           string `json:"name"`           //! 账户用户名
	Password       string `json:"password"`       //! 账户密码
	CreateTime     int64  `json:"createtime"`     //! 账号创建时间
	LastLoginTime  int64  `json:"lastlgtime"`     //! 上次登录时间
	LoginDays      int    `json:"logindays"`      //! 连续登陆天数
	DeviceID       int    `json:"deviceid"`       //! 设备唯一ID
	Forbidden      bool   `json:"forbidden"`      //! 是否禁止登陆
	LastServerID   int    `json:"lastserverid"`   //! 上次选择的服务器ID
	LoginServerIDs []int  `json:"loginserverids"` //! 登录过的游戏服务器ID
}

func Test() {
	var data AccountInfo
	db.Init(Url, PoolLimit)
	db.Find("Account", "AccountsInfo", "_id", 1, &data)
	fmt.Println(data)
}

func TestInsert() {
	db.Init(Url, PoolLimit)

	player1 := Player{
		ID:    1,
		Name:  "小明",
		Level: 43,
		Money: 997653,
		Skill: []int{2003, 4007, 2075}}
	db.Insert("TestDB", "TestTable", &player1)

	player2 := Player{
		ID:    2,
		Name:  "小红",
		Level: 56,
		Money: 8814751,
		Skill: []int{4512, 1024, 4034}}
	db.Insert("TestDB", "TestTable", &player2)

	player3 := Player{
		ID:    3,
		Name:  "二狗",
		Level: 15,
		Money: 250,
		Skill: []int{1015}}
	db.Insert("TestDB", "TestTable", &player3)

	player4 := Player{
		ID:    4,
		Name:  "徐志雷",
		Level: 999,
		Money: 999999999,
		Skill: []int{8971, 9899, 7998, 9999, 9867, 9997}}
	db.Insert("TestDB", "TestTable", &player4)
}

func TestFind() {
	db.Init(Url, PoolLimit)

	player := []Player{}

	//! 条件查找
	db.Find_Conditional("TestDB", "TestTable", "money", ">=", 999999999, &player)
	fmt.Println(player)

	//! 范围查找
	db.Find_Range("TestDB", "TestTable", "level", 10, 100, true, &player)
	fmt.Println(player)

	//! 排序查找
	db.Find_Sort("TestDB", "TestTable", "money", -1, 2, &player)
	fmt.Println(player)
}

func TestLevelup() {
	db.Init(Url, PoolLimit)

	db.IncFieldValue("TestDB", "TestTable", "_id", 4, "level", 1)
}

func TestArrayOperation() {
	db.Init(Url, PoolLimit)

	db.RemoveFromArray("TestDB", "TestTable", "name", "徐志雷", "skill", 7998)
	db.AddToArray("TestDB", "TestTable", "name", "徐志雷", "skill", 25000)
}
