package message

// AckMessage 回复消息
type AckMessage struct {
	Ack  int         `json:"ack"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 正常响应
var Ok = 200

// 不合法消息体
var Invalid = 400

// 失败消息体
var Failed = 500

func AckMessageOk(data interface{}) *AckMessage {
	return &AckMessage{
		Ok,
		"成功",
		data,
	}
}

func AckMessageInvalid(data interface{}) *AckMessage {
	return &AckMessage{
		Invalid,
		"非法请求",
		data,
	}
}

func AckMessageFailed(data interface{}) *AckMessage {
	return &AckMessage{
		Failed,
		"操作失败",
		data,
	}
}

func (this *AckMessage) Invoke() Message {
	return nil
}
