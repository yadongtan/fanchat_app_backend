package znet

import (
	"errors"
	"fantastic_chat/zinx/utils"
	"fantastic_chat/zinx/ziface"
	"fmt"
	"net"
)

// iServer的接口实现
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前Server的消息管理模块
	MsgHandle ziface.IMsgHandle

	//该server的连接管理器
	ConnMgr ziface.IConnManager

	// 该Server创建连接之后自动调用Hook函数 - OnConnStart()
	OnConnStart func(conn ziface.IConnection)
	// 该Server销毁连接之前自动调用的Hook函数-OnConnStOP()
	OnConnStop func(conn ziface.IConnection)
}

// 定义当前客户端连接的所绑定的handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start() {
	// 1.获取一个TCP的Addr
	fmt.Printf("[Zinx] Server Name : %s, listener IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPacketSize : %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err", err)
			return
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listening...")
		var cid uint32
		cid = 0
		// 3.阻塞的等待客户端连接, 处理客户端连接业务(读写)

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//设置最大连接个数的判断, 如果超过最大连接, 那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源, 状态或者一些已经开辟的连接信息进行停止或者回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	// 启动Server的服务功能
	s.Start()
	// 做一些启动服务器之后的额外工作...
	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

// 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}

	return s
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop() ...")
		s.OnConnStop(conn)
	}
}
