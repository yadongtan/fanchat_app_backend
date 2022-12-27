package message

import (
	"fantastic_chat/server/channel"
	"fantastic_chat/server/database"
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
			go c.SendPreviousMsgPublicChat()
			go c.SendPreviousMsgDirectChat()
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
				Ctime:    time.Now().Format("2006-01-02 15:04:05"),
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
			publicTextMsg := <-MsgChan
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

//查询用户最后一次上线的时间, 并将消息发给该用户
func (this *Channel) SendPreviousMsgPublicChat() {
	// 先查询用户最后一次上线时间
	log := &database.UserSigninLog{}
	err := database.GetDB().Where("ttid = ? and `type` = 'offline'", this.TTid).Order("ctime").First(log)
	var t float32
	if log.Ctime != "" && err == nil {
		// 部分消息
		timeStruct, err := time.Parse("2006-01-02 15:04:05", log.Ctime)
		if err != nil {
			fmt.Println("SendPreviousMsgPublicChat() 解析时间错误! err: ", err)
		}
		unixTInt64 := timeStruct.Unix()
		tFloat64, _ := strconv.ParseFloat(strconv.FormatInt(unixTInt64, 10), 32)
		t = float32(tFloat64)
	}
	PublicChatTextMessageArray := GetTextFromRedis(t)
	for _, msg := range PublicChatTextMessageArray {
		if this.TTid == msg.TTid {
			continue
		}
		this.Write(msg)
	}
}

// 查询用户所有历史私聊消息, 并将消息发给该用户
func (this *Channel) SendPreviousMsgDirectChat() {
	// 先查询用户最后一次上线时间 []FriendDirectMessage, error
	msgs, err := GetDirectMsg(0, this.TTid)
	if err != nil {
		return
	}
	// 发送
	for _, msg := range msgs {
		this.Write(&msg)
	}
}
