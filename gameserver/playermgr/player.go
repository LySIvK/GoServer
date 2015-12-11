package PlayerMgr

import (
	"encoding/json"
	"gopath/code.google.com/p/go.net/websocket"
	"loger"
	"net"
	"time"
)

type Player struct {
	PlayerID       int64            //! 玩家ID
	ws             *websocket.Conn  //! WebSocket底层套接字
	sendMsgChannel chan interface{} //! 发送消息频道
	seqID          int              //! 消息顺序号
	PlayerMgr      *PlayerMgr       //! 玩家管理器指针
}

//! 检测消息顺序号
func (self *Player) CheckMsgSeq(msgSeq int) bool {
	if self.seqID != msgSeq {
		loger.Error("checkMsgSeq fail, msgSeq: %d  needSeq: %d", msgSeq, self.seqID)
		return false
	}
	return true
}

//! 处理消息
func (self *Player) ProcessMsg(msg string) {
	//! 顺序号自加
	self.seqID++

	//! 调用消息分拣
	isSuccess := G_Dispatch.DispatchMsg(self, msg)
	if isSuccess == false {
		loger.Error("Dispatch message: %s   has a error.", msg)
	}
}

//! 读取消息协程
func (self *Player) RecvMsg() {
	msg := ""
	defer self.ws.Close()

	for {
		//! 设置1秒超时
		self.ws.SetReadDeadline(time.Now().Add(1 * time.Second))
		err := websocket.Message.Receive(self.ws, &msg)
		if err == nil {
			loger.Debug("PlayerID: %v  recv msg: %s", self.PlayerID, msg)

			//! 交于消息分拣器
			self.ProcessMsg(msg)
			continue
		}

		//! 判断超时
		neterr, ok := err.(net.Error)
		if ok == true && neterr.Timeout() == true {
			time.Sleep(1)
			continue
		}

		//! 其他错误断开连接
		loger.Warn("Disconnect playerID: %v, error: %s", self.PlayerID, err.Error())
		break
	}
}

//! 发送消息协程
func (self *Player) SendMsg() {
	defer self.ws.Close()

	//! 从频道中读取消息
	for msg := range self.sendMsgChannel {
		err := websocket.JSON.Send(self.ws, msg)

		//! 解析消息并输出
		msgText, _ := json.Marshal(msg)
		loger.Debug("Send to id: %v msg: %s", self.PlayerID, msgText)
		if err != nil {
			loger.Error("websocket send msg fail. error: %v", err.Error())
			return
		}
	}
}

//! 发送消息
func (self *Player) Send(msg IMsgHead) {
	msgType, msgAction := msg.GetTypeAndAction()
	msg.FillMsgHead(self.seqID, msgType, msgAction)
	self.sendMsgChannel <- msg
}

//! 读写消息协程运作
func (self *Player) Run() {
	go self.SendMsg()
	self.RecvMsg()
}

//! 生成一个新的Player类
func NewPlayer(ws *websocket.Conn, playerMgr *PlayerMgr) *Player {
	player := new(Player)
	player.ws = ws
	player.PlayerMgr = playerMgr
	player.sendMsgChannel = make(chan interface{}, 1024)
	return player
}
