package getty

import (
	"fantastic_chat/getty/reactor"
	"fantastic_chat/getty/worker"
	"fmt"
)

type GettyEngine struct {
	r *reactor.Reactor //Reactor
	w []*worker.Worker //worker
}

func Bind(ip string, port int) *GettyEngine {

	engine := &GettyEngine{}

	// 创建Reactor
	r := &reactor.Reactor{}
	r.Bind(ip, port)
	engine.r = r
	fmt.Printf("Bind reactor, ip:%s, port:%v\n", ip, port)
	fmt.Printf("reactor info:{%v}\n", r)

	// 创建Worker
	worker_count := 2
	workers := worker.CreateWorker(worker_count)
	engine.w = workers
	fmt.Printf("Create Worker:%v, count:%d\n", workers, worker_count)

	return engine
}

//启动
func (this *GettyEngine) Start() {
	this.r.StartTcp(this.w)
}
