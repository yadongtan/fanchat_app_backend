package handler

import (
	"fmt"
	"net"
)

type Channel struct {
	Ctx *Context //上下文
	Uid int      //用户id
}

// 在线用户集合
var OnlineUserChannelMap map[int]*Channel

// 要添加到在线用户集合的管道
var OnlineUserChannelChan chan *Channel

// 要从在线状态中移除的用户的管道
var OfflineUserChannelChan chan *Channel

func init() {

	OnlineUserChannelChan = make(chan *Channel, 10)
	OfflineUserChannelChan = make(chan *Channel, 10)
	OnlineUserChannelMap = make(map[int]*Channel)
	// 添加在线用户
	go func() {
		for {
			c := <-OnlineUserChannelChan
			OnlineUserChannelMap[c.Uid] = c
			fmt.Printf("用户[uid=%d]已上线\n", c.Uid)
			fmt.Printf("当前在线用户:%v\n", OnlineUserChannelMap)
		}
	}()
	// 移除在线用户
	go func() {
		for {
			c := <-OfflineUserChannelChan
			delete(OnlineUserChannelMap, c.Uid)
			fmt.Printf("当前在线用户:%v\n", OnlineUserChannelMap)
		}
	}()
}

func CreateChannel(conn net.Conn) *Channel {
	ch := &Channel{
		Ctx: &Context{
			conn,
			0,
			&HandlerChain{},
			nil,
		},
		Uid: -1,
	}
	ch.Ctx.Ch = ch
	return ch
}

// 写入
func (this *Channel) Write(b interface{}) {
	this.Ctx.Chain.Write(this.Ctx, b)
}

// 读取
func (this *Channel) Read() (interface{}, error) {
	return this.Ctx.Chain.Read(this.Ctx)
}

// 保持连接并一直读取和处理数据
func (this *Channel) KeepAlive() {
	for {
		_, err := this.Read()
		if err != nil {
			fmt.Printf("用户[uid=%d]已下线\n", this.Uid)
			OfflineUserChannelChan <- this
			return
		}
	}
}

func (this *Channel) AddHandler(h Handler) {
	this.Ctx.Chain.AddHandler(h)
}
