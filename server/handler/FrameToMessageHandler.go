package handler

import (
	"fantastic_chat/server/frame"
	"fantastic_chat/server/message"
	"fmt"
)

type FrameToMessageHandler struct {
}

func (this *FrameToMessageHandler) read(ctx *Context, obj interface{}) (interface{}, error) {

	if obj == nil {
		return nil, nil
	}

	f := obj.(*frame.Frame)
	msg := f.Payload.(message.Message)

	ackMsg := msg.Invoke()

	//判断登录消息
	if f.FrameType == message.SignInFrameType {
		//登录成功
		if ackMsg.(*message.AckMessage).Ack == message.Ok {
			fmt.Printf("用户[%s] 登录成功\n", msg.(*message.SignInMessage).Username)
			ctx.Ch.TTid = msg.(*message.SignInMessage).TTid
			OnlineUserChannelChan <- ctx.Ch
		} else {
			fmt.Printf("用户[%s] 登录失败\n", msg.(*message.SignInMessage).Username)
		}
	}

	fmt.Printf("AckMsg:%v\n", ackMsg)
	ackF := frame.GenerateAckFrame(f, ackMsg)
	fmt.Printf("AckFrame:%v\n", ackF)
	ctx.Chain.triggerNextWriteHandler(ctx, this, ackF)
	return f, nil
}

func (this *FrameToMessageHandler) write(ctx *Context, obj interface{}) interface{} {
	f := frame.GenerateMessageFrame(ctx.Ch.GenerateFrameId(), obj) //将Message包装成帧
	//f是响应结果
	ackF := ctx.Chain.triggerNextWriteHandler(ctx, this, f)

	if ackF == nil {
		return nil
	}

	msg := ackF.(*frame.Frame).Payload.(message.Message)
	if ackF.(*frame.Frame).FrameType != message.AckFrameType {
		fmt.Printf("接收到响应msg:%v \t 但该响应不是Ack类型!\n", msg)
	}
	return msg
}
