package message

// LogoutMessage 登出
type LogoutMessage struct {
	Time string `json:"time"` //下线时间
}

func (this *LogoutMessage) Invoke() Message {
	return nil
}
