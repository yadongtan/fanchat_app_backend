package handler

import (
	"fantastic_chat/server/frame"
	"fantastic_chat/server/message"
	"fmt"
)

var frameId int

func init() {
	frameId = 1
}

type FrameToMessageHandler struct {
}

func (*FrameToMessageHandler) read(ctx *Context, obj interface{}) (interface{}, error) {
	f := obj.(*frame.Frame)
	msg := f.Payload.(message.Message)

	ackMsg := msg.Invoke()

	//判断登录消息
	if f.FrameType == message.SignInFrameType {
		//登录成功
		if ackMsg.(*message.AckMessage).Ack == message.Ok {
			fmt.Printf("用户[%s] 登录成功\n", msg.(*message.SignInMessage).Username)
			ctx.Ch.Uid = msg.(*message.SignInMessage).Uid
			OnlineUserChannelChan <- ctx.Ch
		} else {
			fmt.Printf("用户[%s] 登录失败\n", msg.(*message.SignInMessage).Username)
		}
	}

	fmt.Printf("AckMsg:%v\n", ackMsg)
	ackF := frame.GenerateAckFrame(f, ackMsg)
	ctx.Write(ackF)
	return f, nil
}

func (*FrameToMessageHandler) write(ctx *Context, obj interface{}) interface{} {
	f := frame.GenerateMessageFrame(obj)
	frameId++
	return f
}
