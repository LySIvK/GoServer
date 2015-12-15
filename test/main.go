package main

import (
	"loger"
	//"serverconfig"

	//"loginserver/accountmgr"
)

type TestData struct {
	ID    int
	Name  string
	Level int
	Money int
}

func (self *TestData) GetPathName() string {
	return "test.csv"
}

func (self *TestData) GetName() string {
	return "test"
}

func main() {
	//! 初始化loger
	loger.InitLoger("./log", loger.LogDebug, true, "test")
	loger.Debug("Test Run")

	TestRegister()
	// test := new(TestData)
	// mgr := new(StaticDataMgr)
	// mgr.Init()
	// mgr.Add(test)
	// mgr.Parse()

	//!TestParseCsv()
	//TestLogin()
	//TestCreateRole()
	//TestSendMsg()

	//!TestLogin()   //! Done

	//!TestRegister() //! Done

	//Test()
	//AccountMgr.SayHello()

	//!	serverconfig.Init() //! Done
	//! TestArrayOperation()//! Done
	//! TestLevelup() 		//! Done
	//! TestFind()    		//! Done
	//! TestInsert()  		//! Done
	//! TestSendMsg() 		//! Done
}
