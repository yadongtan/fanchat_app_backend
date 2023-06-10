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

func AckMessageOk(msg string, data interface{}) *AckMessage {

	return &AckMessage{
		Ok,
		msg,
		data,
	}
}

func AckMessageInvalid(msg string, data interface{}) *AckMessage {
	return &AckMessage{
		Invalid,
		msg,
		data,
	}
}

func AckMessageFailed(msg string, data interface{}) *AckMessage {
	return &AckMessage{
		Failed,
		msg,
		data,
	}
}

func (this *AckMessage) Invoke() Message {
	return nil
}
