package handler

import (
	"fantastic_chat/server/utils"
	"fmt"
)

type LengthFieldBasedFrameDecoder struct {
}

func (this *LengthFieldBasedFrameDecoder) read(ctx *Context, obj interface{}) (interface{}, error) {
	// go chain.Read() --> Read()...
	// go chain.Read() ---> ctx.Read() ---> channel.Write()
	currentReadCount := 0           //已经读了的长度
	frameLen := 4                   //要读的长度
	frameLenByte := make([]byte, 8) // 存储长度字节的数组
	// 读取长度
	readLenByte := make([]byte, 4)
	for currentReadCount < frameLen {
		cnt, err := ctx.Conn.Read(readLenByte)
		if err != nil {
			fmt.Printf("[LengthFieldBasedFrameDecoder] Read failed! error: %v\n", err)
			return nil, err
		}
		copy(frameLenByte[currentReadCount:], readLenByte[0:cnt])
		currentReadCount += cnt
	}
	// 更新长度
	frameLen = utils.CastBytesToInt(frameLenByte, 0)
	rawData := make([]byte, frameLen) //存储数据的数组

	copy(rawData, frameLenByte[0:currentReadCount])

	//循环读取
	for currentReadCount < frameLen {
		cnt, err := ctx.Conn.Read(rawData[currentReadCount:])
		if err != nil {
			fmt.Printf("[LengthFieldBasedFrameDecoder] Read failed! error: %v\n", err)
			return nil, err
		}
		currentReadCount += cnt
	}
	fmt.Printf("[LengthFieldBasedFrameDecoder] received data: [ %s ]\n", string(rawData))
	// fmt.Printf("解析到数据帧:%v\n", string(rawData))
	return rawData, nil
}

func (this *LengthFieldBasedFrameDecoder) write(ctx *Context, obj interface{}) interface{} {
	for sent := 0; sent < len(obj.([]byte)); {
		cnt, err := ctx.Conn.Write(obj.([]byte)[sent:])
		if err != nil {
			fmt.Printf("Write to Conn Failed! err: %v\n", err)
			return nil
		}
		sent += cnt
	}
	//这里应该等待消息回应
	return nil
}
