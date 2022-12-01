package message

// 广播消息
type PublicChatTextMessage struct {
	TTid     int    `json:"ttid"`
	Username string `json:"username"`
	Time     string `json:"time"`
	Text     string `json:"text"`
}
