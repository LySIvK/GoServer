package main

import (
	"loger"
	//"serverconfig"

	//"loginserver/accountmgr"
)

func main() {
	//! 初始化loger
	loger.InitLoger("./log", loger.LogDebug, true, "test")
	loger.Debug("Test Run")

	//TestLogin()
	TestCreateRole()
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
