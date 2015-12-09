package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopath/code.google.com/p/go.net/websocket"
	"io"
	"loger"
	"net"
	"net/http"
	"time"
)

var addr = flag.String("0.0.0.0", ":9000", "http service address")

//! 服务器
type Server struct {
	serverID            int          //! 服务器ID
	addClientChannel    chan *Client //! 客户端连入频道
	removeClientChannel chan *Client //! 客户端登出频道
}

//! 客户端处理协程
type Client struct {
	clientID       int              //! 玩家id
	ws             *websocket.Conn  //! 底层套接字
	sendMsgChannel chan interface{} //! 发送消息频道
}

//! 客户端协程
func (self *Client) Run() {
	go self.SendMsg()
	self.RecvMsg()
}

//! 读取消息
func (self *Client) RecvMsg() {
	msg := ""
	defer self.ws.Close()

	for {
		self.ws.SetReadDeadline(time.Now().Add(1 * time.Second)) //! 设置1秒超时
		err := websocket.Message.Receive(self.ws, &msg)
		if nil == err {
			fmt.Println("get clinet msg: ", msg)

			//! TODO: 处理消息
			msg := "World~~"
			self.sendMsgChannel <- msg
			continue
		}

		neterr, ok := err.(net.Error)
		if ok == true && neterr.Timeout() == true {
			time.Sleep(1)
			continue
		}

		//! 其他错误断开连接
		if err == io.EOF {
			fmt.Println("safe quit")
			break
		}
		fmt.Printf("Error: %v \r\n", err)
		break
	}
}

//! 发送消息
func (self *Client) SendMsg() {
	defer self.ws.Close()

	//! 从频道中取出消息
	for msg := range self.sendMsgChannel {
		err := websocket.JSON.Send(self.ws, msg)

		//! 解析消息留下日志
		msgText, _ := json.Marshal(msg)
		fmt.Println("Send to client :", string(msgText))
		if err != nil {
			fmt.Println("Has a error ~!!! ", err)
			break
		}
	}
}

//! 初始化服务器
func (self *Server) InitServer() {
	//! TODO somethings
}

//! 服务器连接回调
func (self *Server) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		client := new(Client)
		client.ws = ws
		client.sendMsgChannel = make(chan interface{}, 1024)

		self.addClientChannel <- client //! 通知服务器一个客户端连上服务

		client.Run()

		self.removeClientChannel <- client //! 断线时通知服务器协程删除此客户端
	}

	return websocket.Handler(connectHandler)
}

//! 服务器监听
func (self *Server) Listen() {
	for {
		select {
		case clinet := <-self.addClientChannel:
			fmt.Println("有一个玩家连入服务器", clinet.ws)
		case clinet := <-self.removeClientChannel:
			fmt.Println("有一个玩家登出服务器", clinet.ws)
		}
	}
}

func NewServer() *Server {
	newServer := new(Server)
	newServer.addClientChannel = make(chan *Client)
	newServer.removeClientChannel = make(chan *Client)
	return newServer
}

func main() {
	loger.InitLoger("./log", loger.LogDebug, true, "test")
	loger.Debug("Debug")
	loger.Error("Error")
	return
	//! Websocket
	flag.Parse()

	serverSingleton := NewServer()
	serverSingleton.InitServer()
	onConnectHandler := serverSingleton.GetConnectHandler()
	http.Handle("/", onConnectHandler)
	go serverSingleton.Listen()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		loger.Fatal("ListenAndServe: ", err)
	}
}
