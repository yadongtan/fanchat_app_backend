package message

// AckMessage 回复消息
type AckMessage struct {
	Ack int
}

// 正常响应
var Ok = 200

// 不合法消息体
var Invalid = 400

// 失败消息体
var Failed = 500

var AckMessageOk *AckMessage
var AckMessageInvalied *AckMessage
var AckMessageFailed *AckMessage

func init() {
	AckMessageOk = &AckMessage{
		Ok,
	}
	AckMessageInvalied = &AckMessage{
		Invalid,
	}
	AckMessageFailed = &AckMessage{
		Failed,
	}
}

func (this *AckMessage) Invoke() Message {
	return nil
}
