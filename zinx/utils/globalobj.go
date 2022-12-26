package utils

import (
	"encoding/json"
	"fantastic_chat/zinx/ziface"
	"io/ioutil"
)

type GlobalObj struct {

	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version        string //当前Zinx的版本
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值

	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //Zinx框架允许用户最多开辟多少个Worker数量
}

/*
	定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

/*
	从zinx.json去加载自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个init方法, 初始化当前的GlobalObject
*/
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 应该尝试从配置文件加载
	//GlobalObject.Reload()
}
