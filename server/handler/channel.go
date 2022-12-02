package handler

import (
	"fantastic_chat/server/channel"
	"fantastic_chat/server/database"
	"fantastic_chat/server/message"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Channel struct {
	Ctx              *Context // 上下文
	TTid             int      // 用户id
	Username         string   // 用户名称
	frameIdIncrement uint32
	wg               sync.WaitGroup
}

// 原子操作版加函数
func (this *Channel) atomicIncrNum(i uint32) int {
	this.wg.Add(1)
	defer this.wg.Done()
	return int(atomic.AddUint32(&i, 1))
}

func (this *Channel) GenerateFrameId() string {
	incIdNum := this.atomicIncrNum(this.frameIdIncrement) //获取一个唯一的自增数字
	frameId := IntToString(this.TTid, 9) + IntToString(incIdNum, 10)
	return frameId
}

func IntToString(i int, numLength int) string {
	//数字右边部分
	right := strconv.Itoa(i)
	//左边差多少位补多少个0
	count := numLength - len(right)

	left := ""
	for j := 0; j < count; j++ {
		left += "0"
	}
	return left + right

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
			OnlineUserChannelMap[c.TTid] = c
			channel.Cs.OnlinePersonCount = len(OnlineUserChannelMap)
			fmt.Printf("用户[TTid=%d]已上线\n", c.TTid)
			fmt.Printf("当前在线用户:%v\n", OnlineUserChannelMap)

		}
	}()
	// 移除在线用户
	go func() {
		for {
			c := <-OfflineUserChannelChan
			delete(OnlineUserChannelMap, c.TTid)

			// 添加离线日志
			ip := strings.Split(c.Ctx.Conn.RemoteAddr().String(), ":")[0]
			ipDetails := database.GetIpDetails(ip)

			signinLog := &database.UserSigninLog{
				TTid:     c.TTid,
				Type:     "Offline",
				Ctime:    time.Now().Format("2002-01-01 01:01:01"),
				Ip:       ip,
				Province: ipDetails.Province,
				City:     ipDetails.City,
				Region:   ipDetails.Region,
				Addr:     ipDetails.Addr,
			}

			database.GetDB().Create(signinLog)
			fmt.Printf("当前在线用户:%v\n", OnlineUserChannelMap)
		}
	}()
	// 写出所有公有消息
	go func() {
		for {
			publicTextMsg := <-message.MsgChan
			for ttid, ch := range OnlineUserChannelMap {
				if ttid == publicTextMsg.TTid {
					continue
				}

				go ch.Write(publicTextMsg)
			}
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
		TTid: -1,
		wg:   sync.WaitGroup{},
	}
	ch.Ctx.Ch = ch
	return ch
}

// 写入
func (this *Channel) Write(b interface{}) interface{} {
	return this.Ctx.Chain.triggerNextWriteHandler(this.Ctx, nil, b)
}

// 读取
func (this *Channel) Read() (interface{}, error) {
	return this.Ctx.Chain.Read(this.Ctx)
}

// 保持连接并一直读取和处理数据
func (this *Channel) KeepAlive() {
	//测试
	//go func() {
	//	msg := &message.SignInMessage{
	//		2209931449,
	//		"yadong",
	//		"yadong123456",
	//	}
	//	// _, err = conn.Write([]byte(message))
	//	time.Sleep(15 * time.Second)
	//	ackMsg := this.Write(msg)
	//	fmt.Printf("接收到AckMsg : %v\n", ackMsg)
	//}()

	for {
		_, err := this.Read()
		if err != nil {
			fmt.Printf("用户[TTid=%d]已下线\n", this.TTid)
			OfflineUserChannelChan <- this
			return
		}

	}
}

func (this *Channel) AddHandler(h Handler) {
	this.Ctx.Chain.AddHandler(h)
}
