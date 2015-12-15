package module

type ChatMgr struct {
	moduleMgr *ModuleMgr //! 模块管理器指针
}

//! 注册消息
func (self *ChatMgr) registerMsg() {
}

//! 初始化管理器
func (self *ChatMgr) Init(moduleMgr *ModuleMgr) {
	G_Dispatch.AddMsgRegistryToMap(new(Msg_Chat_Private)) //! 用户私聊消息
}

//! 生成聊天管理器
func NewChatMgr(moduleMgr *ModuleMgr) *ChatMgr {
	chat := new(ChatMgr)
	chat.Init(moduleMgr)
	return chat
}
