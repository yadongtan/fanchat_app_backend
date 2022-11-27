package handler

import (
	"fmt"
	"net"
)

// 一个用户对应channel

type Context struct {
	Conn         net.Conn
	HandlerIndex int
	Chain        *HandlerChain
	Ch           *Channel
}

func (this *Context) Read() (interface{}, error) {
	return this.Chain.readFromIndex(this, nil, this.HandlerIndex)
}

func (this *Context) Write(b interface{}) {
	this.Chain.writeFromIndex(this, b, this.HandlerIndex)
}

type HandlerChain struct {
	Has []Handler
}

type Handler interface {
	read(ctx *Context, obj interface{}) (interface{}, error)
	write(ctx *Context, obj interface{}) interface{}
}

//从头读
func (this *HandlerChain) Read(ctx *Context) (interface{}, error) {
	return this.readFromIndex(ctx, nil, 0)
}

//从尾写
func (this *HandlerChain) Write(ctx *Context, obj interface{}) {
	this.writeFromIndex(ctx, obj, len(this.Has))
}

//从指定索引开始读
func (this *HandlerChain) readFromIndex(ctx *Context, obj interface{}, index int) (interface{}, error) {

	// 执行
	for i := index; i < len(this.Has); i++ {
		ctx.HandlerIndex = i
		anew, err := this.Has[i].read(ctx, obj)
		obj = anew
		fmt.Printf("[%T] Read: %s\n", this.Has[i], anew)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil

}

//从指定索引开始写
func (this *HandlerChain) writeFromIndex(ctx *Context, obj interface{}, index int) {
	for i := index - 1; i >= 0; i-- {
		ctx.HandlerIndex = i
		obj = this.Has[i].write(ctx, obj)
	}
	_, err := ctx.Conn.Write(obj.([]byte))
	if err != nil {
		fmt.Printf("write to Conn Failed! err : %v\n", err)
		return
	}
}

func (this *HandlerChain) AddHandler(h Handler) {
	this.Has = append(this.Has, h)
}
