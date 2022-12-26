package znet

import (
	"errors"
	"fantastic_chat/zinx/utils"
	"fantastic_chat/zinx/ziface"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {

	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 当前连接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc

	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool

	// 无缓存的管道, 用于读,写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID 和 对应的处理业务API关系
	MsgHandle ziface.IMsgHandle

	//连接属性集合
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandle,
		msgChan:   make(chan []byte),
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中, 最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		// 创建一个拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的MsgHead 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		// 拆包, 得到msgID 和 msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据dataLen, 再次读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
			// 已经开启了工作池机制, 将消息发送给Worker工作池处理即可
		} else {
			// 在执行注册的路由方法
			// 从路由中,找到注册绑定的Conn对应的router调用
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}
}

// 写消息的Goroutine, 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	// 不断阻塞的等待channel的消息, 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出, 此时Writer也要退出
			return
		}
	}
}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// 启动从当前连接的读数据业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的, 创建连接之后需要调用的业务处理, 执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的 销毁连接之前 需要执行的业务Hook汉化
	c.TcpServer.CallOnConnStart(c)

	// 关闭socket连接
	c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态, IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据. 将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加一个连接属性
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//  删除属性
	delete(c.property, key)
}
