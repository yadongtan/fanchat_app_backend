package message

import (
	"fmt"
)

type Message interface {
	Invoke() Message
}

var AckFrameType = 100
var SignInFrameType = 101
var LogoutFrameType = 102

var GetOnlinePersonCountMessageFrameType = 200
var PublicChatImgMessageFrameType = 201
var PublicChatTextMessageFrameType = 202
var PublicChatVideoMessageType = 203

var SignUpMessageFrameType = 301

func GetMessageByType(messageType int) interface{} {
	switch messageType {
	case SignInFrameType:
		return &SignInMessage{}
	case LogoutFrameType:
		return &LogoutMessage{}
	case AckFrameType:
		return &AckMessage{}
	}

	switch messageType {
	case GetOnlinePersonCountMessageFrameType:
		return &GetOnlinePersonCountMessage{}
	case PublicChatImgMessageFrameType:
		return &PublicChatImgMessage{}
	case PublicChatTextMessageFrameType:
		return &PublicChatTextMessage{}
	case PublicChatVideoMessageType:
		return &PublicChatVideoMessage{}
	}

	switch messageType {
	case SignUpMessageFrameType:
		return &SignUpMessage{}
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

	switch msg.(type) {
	case *GetOnlinePersonCountMessage:
		return GetOnlinePersonCountMessageFrameType
	case *PublicChatImgMessage:
		return PublicChatImgMessageFrameType
	case *PublicChatTextMessage:
		return PublicChatTextMessageFrameType
	case *PublicChatVideoMessage:
		return PublicChatVideoMessageType
	}

	switch msg.(type) {
	case *SignUpMessage:
		return SignUpMessageFrameType

	}
	return 0
}
