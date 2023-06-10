package message

import (
	"encoding/json"
	"fantastic_chat/server/myredis"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 广播消息
type PublicChatTextMessage struct {
	TTid     int    `json:"ttid"`
	Username string `json:"username"`
	Time     int64  `json:"time"`
	Text     string `json:"text"`
}

var publicChatMsgInRedisZsetName = "chat:public"

var MsgChan chan *PublicChatTextMessage

func init() {
	MsgChan = make(chan *PublicChatTextMessage, 1024)
}

func (this *PublicChatTextMessage) Invoke() Message {

	// 同时存到Redis
	_ = DurationTextToRedis(this)
	// 转发给其他人
	MsgChan <- this
	return AckMessageOk("Ok", nil)
}

// 持久化此消息到Redis
func DurationTextToRedis(msg *PublicChatTextMessage) error {
	j, _ := json.Marshal(msg)
	_, err := myredis.Client.Get().Do("zadd", publicChatMsgInRedisZsetName, msg.Time, j)
	if err != nil {
		fmt.Println("DurationToRedis() Failed!!! err: ", err)
		return err
	}
	return nil
}

// 获取从指定时间之后的消息
func GetTextFromRedis(fromTime float32) []*PublicChatTextMessage {
	reply, err := redis.Strings(myredis.Client.Get().Do("zrangebyscore", publicChatMsgInRedisZsetName, fromTime, "+inf"))

	if err != nil {
		fmt.Println("GetTextFromRedis() Failed!!! err: ", err)
	}
	pctmMap := make([]*PublicChatTextMessage, len(reply))
	for index, str := range reply {
		pctm := &PublicChatTextMessage{}
		err := json.Unmarshal(([]byte)(str), pctm)
		if err != nil {
			fmt.Println("GetTextFromRedis() Failed!!! err: ", err)
		}
		pctmMap[index] = pctm

	}
	fmt.Printf("reply : %v\n", reply)
	return pctmMap
}
