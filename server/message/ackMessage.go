package message

// AckMessage 回复消息
type AckMessage struct {
	Ack int    `json:"ack"`
	Msg string `json:"msg"`
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
		"成功",
	}
	AckMessageInvalied = &AckMessage{
		Invalid,
		"非法请求",
	}
	AckMessageFailed = &AckMessage{
		Failed,
		"操作失败",
	}
}

func (this *AckMessage) Invoke() Message {
	return nil
}
