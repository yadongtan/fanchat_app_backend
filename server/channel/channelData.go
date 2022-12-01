package channel

type ChannelStatus struct {
	OnlinePersonCount int
}

var Cs ChannelStatus

func init() {
	Cs = ChannelStatus{0}
}
