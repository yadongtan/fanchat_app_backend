package main

import (
	"fantastic_chat/server/message"
	"fmt"
	"net"
	"time"
)

func main() {
	// startClient("yadong.icu", 8081)

	//f := frame.ResolveFrame(bytes)
	//fmt.Printf("加密后: %s\n", bytes)
	//fmt.Printf("解密后:%s\n", f)

	startClient("127.0.0.1", 8081)
}

func startClient(ip string, port int) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	go readData(conn)
	go writeData(conn)
	time.Sleep(3600 * 24 * time.Second)
}

//读取从服务器来的数据
func readData(conn net.Conn) {
	for {
		data := make([]byte, 1024)
		cnt, err := conn.Read(data)
		if err != nil {
			fmt.Println("客户端读取服务器数据时发生错误")
			time.Sleep(10 * time.Second)
			continue
		}
		str := string(data[0:cnt])
		f := message.ResolveFrame([]byte(str))
		//读取到数据
		fmt.Printf("接收到服务器响应: %v\n", f)
	}
}

// 向服务器发送数据
func writeData(conn net.Conn) {
	for {
		// var message string
		//_, err := fmt.Scan(&message)
		//if err != nil {
		//	fmt.Println("输入错误")
		//	continue
		//}
		msg := &message.SignInMessage{
			2209931449,
			"yadong",
			"yadong123456",
		}
		bytes := message.GenerateFrameBytesDefault(1, msg)

		fmt.Println(string(bytes))
		// _, err = conn.Write([]byte(message))
		_, err := conn.Write(bytes)
		if err != nil {
			fmt.Println("网络发生错误, 请稍后重试")
			time.Sleep(10 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)
	}
}
