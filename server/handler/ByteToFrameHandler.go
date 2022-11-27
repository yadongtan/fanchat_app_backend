package handler

import (
	"fantastic_chat/server/frame"
)

type ByteToFrameHandler struct {
}

func (*ByteToFrameHandler) read(ctx *Context, obj interface{}) (interface{}, error) {
	rawData := obj.([]byte)
	f := frame.ResolveFrame(rawData)
	return f, nil
}

func (*ByteToFrameHandler) write(ctx *Context, obj interface{}) interface{} {
	// Frame ---> []Byte
	return frame.GenerateFrameByesWithFrame(obj.(*frame.Frame))
}
