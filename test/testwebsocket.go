package main

import (
	"fmt"
	"gopath/code.google.com/p/go.net/websocket"
)

var url = "ws://192.168.1.102:9000/"
var origin = "http://192.168.1.102/"

//! websocket 连接测试
func TestSendMsg() {
	//! 使用默认配置连接服务端
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println("connect server fail!")
		return
	}

	err = websocket.JSON.Send(ws, "Hello~~")
	if err != nil {
		fmt.Println("send to server fail!")
		return
	}

	fmt.Println("send msg: ", "Hello~~")

	msg := ""
	websocket.JSON.Receive(ws, &msg)

	fmt.Println("recv msg: ", msg)

	fmt.Println("Success!")
	defer ws.Close()
}
