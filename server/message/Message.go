package message

import (
	"fmt"
)

type Message interface {
	Invoke() Message
}

var AckFrameType = 100    //应答
var SignInFrameType = 101 //登录
var LogoutFrameType = 102 //注销

var GetOnlinePersonCountMessageFrameType = 200 //获取在线人数
var PublicChatImgMessageFrameType = 201        //公共频道图片
var PublicChatTextMessageFrameType = 202       //公平频道文字
var PublicChatVideoMessageType = 203           //公共频道视频

var SignUpMessageFrameType = 301 //注册

var UpdateUserAccountMessageFrameType = 401 //更新用户信息
var AddFriendMessageFrameType = 402         //添加朋友
var GetFriendsListMessageFrameType = 403    //查询好友列表
var DirectFriendsMessageFrameType = 404     //私聊消息

var CreateChatMessageFrameType = 501 //创建OpenAI的Chat(Text)聊天机器人
var GetAIFriendsListFrameType = 502

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

	switch messageType {
	case UpdateUserAccountMessageFrameType:
		return &UpdateUserAccountMessage{}
	case AddFriendMessageFrameType:
		return &AddFriendMessage{}
	case GetFriendsListMessageFrameType:
		return &GetFriendsListMessage{}
	case DirectFriendsMessageFrameType:
		return &FriendDirectMessage{}
	}

	switch messageType {
	case CreateChatMessageFrameType:
		return &CreateChatMessage{}
	case GetAIFriendsListFrameType:
		return &GetAIFriendsListMessage{}
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

	switch msg.(type) {
	case *UpdateUserAccountMessage:
		return UpdateUserAccountMessageFrameType
	case *AddFriendMessage:
		return AddFriendMessageFrameType
	case *GetFriendsListMessage:
		return GetFriendsListMessageFrameType
	case *FriendDirectMessage:
		return DirectFriendsMessageFrameType
	}

	switch msg.(type) {
	case *CreateChatMessage:
		return CreateChatMessageFrameType
	case *GetAIFriendsListMessage:
		return GetAIFriendsListFrameType
	}
	return 0
}
