package message

import (
	"fmt"
)

// frameId

type ByteToFrameHandler struct {
	WaitForAckMessageChanMap map[string]chan *Frame
}

func (this *ByteToFrameHandler) read(ctx *Context, obj interface{}) (interface{}, error) {
	rawData := obj.([]byte)
	f := ResolveFrame(rawData)
	fmt.Printf("接收到Frame: %v\n", f)
	//尝试获取chan, 以便回复消息
	ackMessageChan := this.WaitForAckMessageChanMap[f.FrameId]
	if ackMessageChan != nil {
		ackMessageChan <- f //响应消息, 交给下面的write函数
		return nil, nil
	}

	return f, nil
}

func (this *ByteToFrameHandler) write(ctx *Context, obj interface{}) interface{} {

	f := obj.(*Frame)
	frameId := obj.(*Frame).FrameId //frameId
	// Frame ---> []Byte
	frameBytes := GenerateFrameByesWithFrame(obj.(*Frame))

	// 让下一个handler写出, 写出完了以后呢, 创建一个等待响应的chan
	ctx.Chain.triggerNextWriteHandler(ctx, this, frameBytes)

	// 如果是发送Ack, 那么就不用管
	if f.FrameType == AckFrameType {
		return nil
	}

	//等待回复
	waitChan := make(chan *Frame)
	defer close(waitChan)

	this.WaitForAckMessageChanMap[frameId] = waitChan

	ackF := <-waitChan //拿到这个响应的frame

	//删除接收响应的chan
	delete(this.WaitForAckMessageChanMap, frameId)

	fmt.Printf("接收响应frame:%v\n", ackF)

	return ackF

}
