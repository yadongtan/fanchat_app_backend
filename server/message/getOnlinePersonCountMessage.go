package message

import (
	"fantastic_chat/server/channel"
)

type GetOnlinePersonCountMessage struct {
}

func (this *GetOnlinePersonCountMessage) Invoke() Message {
	count := channel.Cs.OnlinePersonCount
	ackMsg := AckMessageOk(count)
	return ackMsg
}
