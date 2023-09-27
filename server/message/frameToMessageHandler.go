package message

import (
	"fmt"
	"strings"
)

type FrameToMessageHandler struct {
}

func (this *FrameToMessageHandler) read(ctx *Context, obj interface{}) (interface{}, error) {

	if obj == nil {
		return nil, nil
	}

	f := obj.(*Frame)
	msg := f.Payload.(Message)
	var ackMsg Message
	//判断登录消息
	if f.FrameType == SignInFrameType {
		msg.(*SignInMessage).Ip = strings.Split(ctx.Conn.RemoteAddr().String(), ":")[0]
		ackMsg = msg.Invoke()
		//登录成功
		if ackMsg.(*AckMessage).Ack == Ok {
			fmt.Printf("用户[%s] 登录成功\n", msg.(*SignInMessage).Username)
			ctx.Ch.TTid = msg.(*SignInMessage).TTid
			ctx.Ch.Username = ackMsg.(*AckMessage).Data.(string)
			OnlineUserChannelChan <- ctx.Ch
		} else {
			fmt.Printf("用户[%s] 登录失败\n", msg.(*SignInMessage).Username)
		}
	} else {
		ackMsg = msg.Invoke()
	}
	fmt.Printf("AckMsg:%v\n", ackMsg)
	ackF := GenerateAckFrame(f, ackMsg)
	fmt.Printf("AckFrame:%v\n", ackF)
	ctx.Chain.triggerNextWriteHandler(ctx, this, ackF)
	return f, nil
}

func (this *FrameToMessageHandler) write(ctx *Context, obj interface{}) interface{} {
	f := GenerateMessageFrame(ctx.Ch.GenerateFrameId(), obj) //将Message包装成帧
	//f是响应结果
	fmt.Printf("发送消息Frame: %v\n", f)
	ackF := ctx.Chain.triggerNextWriteHandler(ctx, this, f)

	if ackF == nil {
		return nil
	}
	msg := ackF.(*Frame).Payload.(Message)
	if ackF.(*Frame).FrameType != AckFrameType {
		fmt.Printf("接收到响应msg:%v \t 但该响应不是Ack类型!\n", msg)
	}
	return msg
}
