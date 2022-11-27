package message

import "fmt"

type Message interface {
	Invoke() Message
}

var AckFrameType = 100
var SignInFrameType = 101
var LogoutFrameType = 102

func GetMessageByType(messageType int) interface{} {
	switch messageType {
	case SignInFrameType:
		return &SignInMessage{}
	case LogoutFrameType:
		return &LogoutMessage{}
	case AckFrameType:
		return &AckMessage{}
	}
	return nil
}

func GetMessageTypeByInterface(msg interface{}) int {
	fmt.Printf("%T\n", msg)
	switch msg.(type) {
	case *SignInMessage:
		return SignInFrameType
	case *LogoutMessage:
		return LogoutFrameType
	case *AckMessage:
		return AckFrameType
	}
	return 0
}
