package worker

import (
	"fmt"
	"net"
)

type Worker struct {
	task     chan func()
	connChan chan net.Conn
	conns    []net.Conn //连接
	hasTask  chan bool
}

//向worker提交任务
func (w *Worker) Submit(task func()) {
	w.task <- task
	w.hasTask <- true
}

//启动worker
func (w *Worker) run() {
	go w.runTask()

}

//暴露给外部添加连接的函数
func (w *Worker) AddNewConn(conn net.Conn) {
	w.connChan <- conn
	w.hasTask <- true
}

func (w *Worker) runTask() {
	fmt.Printf("runTask()...\n")
	for {
		select {
		case f := <-w.task:
			fmt.Printf("处理普通任务\n")
			f()
		case conn := <-w.connChan:
			fmt.Println("处理连接事件\n")
			w.conns = append(w.conns, conn)
			w.executeEvent(conn)
		default:
			fmt.Printf("无任务, 等待唤醒\n")
			if hasTask := <-w.hasTask; hasTask {
				continue
			}
		}
	}
}

// 创建指定数量的worker
func CreateWorker(n int) []*Worker {
	workers := make([]*Worker, n)
	for i := 0; i < n; i++ {
		newWorker := &Worker{}
		workers[i] = newWorker
		newWorker.run()
	}
	return workers
}

//为连接创建读写事件
func (w *Worker) executeEvent(conn net.Conn) {
	go read(conn)
}

//创建读事件
func read(conn net.Conn) {
	for {
		rawBytes := make([]byte, 1024)
		cnt, nil := conn.Read(rawBytes)
		if nil != nil {
			fmt.Printf("read failed in connection: %v\n", conn)
		}
		datas := rawBytes[0:cnt]
		fmt.Printf("connection[%v]received: %s", conn, string(datas))
	}
}
