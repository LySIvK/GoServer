package PlayerMgr

import (
	"Runtime/debug"
	"encoding/json"
	"loger"
	"time"
	"tool"
)

//! 消息头
type MsgHead struct {
	SeqID      int    `json:"seq"`    //! 顺序标志
	MsgType    string `json:"type"`   //! 消息类型
	Action     string `json:"action"` //! 消息行为
	CreateTime int64  `json:"time"`   //! 生成时间
}

type Head struct {
	MsgHead `json:"head"`
}

//! 填充消息头
func (self *MsgHead) FillMsgHead(seqID int, msgType string, msgAction string) {
	if self.SeqID <= 0 { //! 顺序号错误,直接放弃
		self.SeqID = seqID
		self.MsgType = msgType
		self.Action = msgAction
		self.CreateTime = time.Now().Unix()
	}
}

//! 获取消息注册键值
func (self *MsgHead) GetMsgKey() string {
	itemKey := self.MsgType + "-" + self.Action
	return itemKey
}

//! 定义消息接口
type IMsgHead interface {
	FillMsgHead(seqID int, msgType string, msgAction string) //! 填写消息头
	GetTypeAndAction() (string, string)                      //! 获取消息类型与行为
	ProcessAction(player *Player) bool                       //! 处理消息函数
}

type MsgRegistryMap map[string]IMsgHead

//! 消息分拣器
type MsgDispatch struct {
	msgRegistryMap MsgRegistryMap
}

//! 初始化消息分拣器
var G_Dispatch MsgDispatch

func (self *MsgDispatch) Init() {
	self.msgRegistryMap = make(MsgRegistryMap)
}

//! 增加消息子类型
func (self *MsgDispatch) AddMsgRegistryToMap(msg IMsgHead) {
	msgType, msgAction := msg.GetTypeAndAction()
	msgKey := msgType + "-" + msgAction
	_, ok := self.msgRegistryMap[msgKey]
	if ok == true {
		//! 有重复的消息
		loger.Fatal("MsgDispatch addMsgRegistryToMap duplicate: %s", msgKey)
		return
	}

	self.msgRegistryMap[msgKey] = msg
}

//! 获取消息头
func (self *MsgDispatch) GetMsgHead(msg *string) *Head {
	headMsg := new(Head)
	err := json.Unmarshal([]byte(*msg), headMsg)
	if err != nil {
		loger.Warn("Get msg head fail. error: %s", err.Error())
		return nil
	}

	loger.Debug("Head: seq = %d", headMsg.SeqID)
	return headMsg
}

//! 消息分拣
func (self *MsgDispatch) DispatchMsg(player *Player, msg string) bool {
	headMsg := self.GetMsgHead(&msg)
	if headMsg == nil {
		return false
	}

	if player.CheckMsgSeq(headMsg.SeqID) == false {
		return false
	}

	//! 消息处理
	result, skip := self.DealWithMsg(player, headMsg, &msg)
	if skip == false {
		return result
	}

	return false
}

//! 消息处理
func (self *MsgDispatch) DealWithMsg(player *Player, head *Head, msg *string) (bool, bool) {
	defer func() {
		x := recover()
		if x != nil {
			//! 如果错误,则输出栈信息
			loger.Error("%v\r\n%s", x, debug.Stack())
		}
	}()

	itemKey := head.GetMsgKey()
	msgRegistry := self.msgRegistryMap[itemKey]
	if msgRegistry == nil {
		return false, true //! 失败, 跳过该消息
	}

	//! 克隆一个空值对象
	newObj := tool.CloneType(msgRegistry)

	//! 空值对象转型操作
	newMsg := newObj.(IMsgHead)

	//! 解析到对象
	err := json.Unmarshal([]byte(*msg), newMsg)
	if err != nil {
		loger.Error("DealWithMsg unmarshal fail. error: %s", err.Error())
		return false, false //! 直接失败跳出
	}

	//! 消息处理
	result := newMsg.ProcessAction(player)

	return result, false
}

//! 初始化消息处理器列表
func (self *MsgDispatch) initMsgHandlerList() {

}
