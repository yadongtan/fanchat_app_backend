package message

// LogoutMessage 登出
type LogoutMessage struct {
	Time float32 `json:"time"` //下线时间
}

func (this *LogoutMessage) Invoke() Message {
	return nil
}
