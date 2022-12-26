package ziface

type IServer interface {
	Start()

	Stop()

	Server()

	//路由功能: 给当前的服务注册一个路由方法, 供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)

	GetConnMgr() IConnManager

	// 注册OnConnStart钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	// 注册OnConnStop钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// 注册OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)
	// 调用OnConnStop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
