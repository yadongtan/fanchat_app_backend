package reactor

import (
	"fantastic_chat/getty/worker"
	"fmt"
	"net"
)

type Reactor struct {
	ip   string
	port int
}

func (this *Reactor) Bind(ip string, port int) {
	this.ip = ip
	this.port = port
}

func (this *Reactor) BindLocalhost(port int) {
	this.ip = ""
	this.port = port
}

func (this *Reactor) StartTcp(workers []*worker.Worker) {
	address := fmt.Sprintf("%s:%d", this.ip, this.port)
	fmt.Println(address)
	listner, err := net.Listen("tcp", address) //监听
	if err != nil {
		fmt.Printf("Reactor start Listen failed, err:%v\n", err)
		return
	}
	worker_round_robind_index := 0
	fmt.Printf("Start Getty Server successful on Port(s):%d\n", this.port)
	for {
		conn, err := listner.Accept()
		fmt.Printf("接收到客户端连接:%v", conn)
		if err != nil {
			fmt.Printf("Accept connection failed, err:%v\n", err)
		}
		//提交给某一个worker
		workers[worker_round_robind_index].AddNewConn(conn)
	}
}
