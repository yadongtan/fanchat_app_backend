package main

import (
	"container/list"
	"fantastic_chat/server/channel"
	"fantastic_chat/server/message"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type fantasticChatServer struct {
}

var mutexOnlineUserMap sync.Mutex

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ip := ""
	port := 8081

	startServer(ip, port)

}

func startServer(ip string, port int) {
	connList := list.New()
	address := fmt.Sprintf("%s:%d", ip, port)
	// 监听端口
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	//接收客户端请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Listen err: ", err)
			continue
		}

		ch := message.CreateChannel(conn)
		ch.AddHandler(&message.LengthFieldBasedFrameDecoder{})
		ch.AddHandler(&message.ByteToFrameHandler{
			make(map[string]chan *message.Frame),
		})
		ch.AddHandler(&message.FrameToMessageHandler{})

		go ch.KeepAlive() // 持续读取并处理数据

		connList.PushFront(&conn)
		fmt.Println("当前已连接的客户端:", connList)
		//// 分发客户端上线消息
		//go sendOnlineMessage(conn, connList)
		//// 读取并处理客户端数据
		//go readHandler(conn, connList)
	}
}

//添加在线用户
func (this *fantasticChatServer) addOnlineUser(uid int, conn net.Conn) {
	mutexOnlineUserMap.Lock()
	channel.Cs.OnlineUserMap[uid] = conn
	mutexOnlineUserMap.Unlock()
}

//用户上线通知
func sendOnlineMessage(conn net.Conn, connList *list.List) {
	addr := conn.RemoteAddr()
	online := fmt.Sprintf("%s 已上线", addr)
	fmt.Println(online)
	sendMessageToAll(online, connList)
}

//读取数据转发给其他人
func readHandler(conn net.Conn, connList *list.List) {
	for {
		rawBytes := make([]byte, 1024)

		//         4                 4                 4                     4                      4                 4                      _
		//  -----------------------------------------------------------------------------------------------------------------------------------------
		//  | frameLen(帧长度) | version(版本号) | frameId(帧id) | frameType(帧类型) | serializeType(序列化类型) | encryptType(加密类型) | payload(消息内容) |
		//  -----------------------------------------------------------------------------------------------------------------------------------------
		cnt, err := conn.Read(rawBytes)
		if err != nil {
			return
		}
		frame := string(rawBytes[0:cnt])
		details := fmt.Sprintf("用户[%s]: %s", conn.RemoteAddr(), frame)
		sendMessageToAll(details, connList)
	}
}

//转发给所有在线用户
func sendMessageToAll(message string, connList *list.List) {
	for i := connList.Front(); i != nil; i = i.Next() {
		_, err := (*i.Value.(*net.Conn)).Write([]byte(message))
		if err != nil {
			continue
		}
	}
}
