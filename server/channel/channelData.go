package channel

import "net"

type ChannelStatus struct {
	OnlinePersonCount int
	OnlineUserMap     map[int]net.Conn //用户id与conn
}

var Cs ChannelStatus

func init() {
	Cs = ChannelStatus{0, make(map[int]net.Conn)}
}
