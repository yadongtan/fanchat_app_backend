package message

// 广播消息
type PublicChatVideoMessage struct {
	TTid     int    `json:"ttid"`
	Username string `json:"username"`
}

func (this *PublicChatVideoMessage) Invoke() Message {
	return AckMessageOk("Ok", nil)
}
