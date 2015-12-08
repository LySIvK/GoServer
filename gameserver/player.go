package main

import (
	"encoding/json"
	"gopath/code.google.com/p/go.net/websocket"
	"loger"
	"net"
	"time"
)

type Player struct {
	playerID       int64            //! 玩家ID
	ws             *websocket.Conn  //! WebSocket底层套接字
	sendMsgChannel chan interface{} //! 发送消息频道
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
			loger.Debug("PlayerID: %v  recv msg: %s", self.playerID, msg)

			//! TODO: 消息处理
			continue
		}

		//! 判断超时
		neterr, ok := err.(net.Error)
		if ok == true && neterr.Timeout() == true {
			time.Sleep(1)
			continue
		}

		//! 其他错误断开连接
		loger.Warn("Disconnect playerID: %v, error: %s", self.playerID, err.Error())
		break
	}
}

//! 发送消息协程
func (self *Player) SendMsg() {
	defer self.ws.Close()

	//! 从频道中读取消息
	for msg := range self.sendMsgChannel {
		err := websocket.Message.Send(self.ws, msg)

		//! 解析消息并输出
		msgText, _ := json.Marshal(msg)
		loger.Debug("Send to id: %v  msg: ", self.playerID, msgText)
		if err != nil {
			loger.Error("websocket send msg fail. error: %v", err.Error())
			return
		}
	}
}

//! 发送消息
func (self *Player) Send(msg string) {
	self.sendMsgChannel <- msg
}

//! 读写消息协程运作
func (self *Player) Run() {
	go self.SendMsg()
	self.RecvMsg()
}
