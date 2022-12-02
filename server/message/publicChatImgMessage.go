package message

// 广播消息
type PublicChatImgMessage struct {
	TTid     int    `json:"ttid"`
	Username string `json:"username"`
}

func (this *PublicChatImgMessage) Invoke() Message {
	return AckMessageOk(nil)
}
