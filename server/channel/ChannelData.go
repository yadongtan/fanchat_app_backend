package channel

import (
	"net"
	"sync/atomic"
)

type ChannelStatus struct {
	OnlinePersonCount int32
	OnlineUserMap     map[int]net.Conn //用户id与conn
}

var Cs = ChannelStatus{0, make(map[int]net.Conn)}

func IncrOnlinePersonCount() {
	atomic.AddInt32(&Cs.OnlinePersonCount, 1)
}

func DecrOnlinePersonCount() {
	atomic.AddInt32(&Cs.OnlinePersonCount, -1)
}
