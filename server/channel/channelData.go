package channel

import (
	"net"
	"sync"
	"sync/atomic"
)

type ChannelStatus struct {
	OnlinePersonCount int32
	OnlineUserMap     map[int]net.Conn //用户id与conn
}

var wg sync.WaitGroup
var Cs ChannelStatus

func init() {
	Cs = ChannelStatus{0, make(map[int]net.Conn)}
	wg = sync.WaitGroup{}
}

func IncrOnlinePersonCount() {
	AtomicOperateNum(&Cs.OnlinePersonCount, 1)
}

func DecrOnlinePersonCount() {
	AtomicOperateNum(&Cs.OnlinePersonCount, -1)
}

// 原子操作版加函数
func AtomicOperateNum(i *int32, step int32) int {
	wg.Add(1)
	defer wg.Done()
	return int(atomic.AddInt32(i, step))
}
