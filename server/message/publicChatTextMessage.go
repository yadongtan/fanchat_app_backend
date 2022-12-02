package message

import (
	"encoding/json"
	"fantastic_chat/server/redis"
	"fmt"
)

// 广播消息
type PublicChatTextMessage struct {
	TTid     int     `json:"ttid"`
	Username string  `json:"username"`
	Time     float32 `json:"time"`
	Text     string  `json:"text"`
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
	_, err := redis.Client.Get().Do("zadd", publicChatMsgInRedisZsetName, msg.Time, j)
	if err != nil {
		fmt.Println("DurationToRedis() Failed!!! err: ", err)
		return err
	}
	return nil
}

// 获取从指定时间之后的消息
func GetTextFromRedis(fromTime string) []PublicChatTextMessage {
	reply, err := redis.Client.Get().Do("zrangebyscore", publicChatMsgInRedisZsetName, fromTime, "+inf")
	if err != nil {
		fmt.Println("GetTextFromRedis() Failed!!! err: ", err)
	}
	fmt.Println("GetTextFromRedis() reply:", reply.(string))
	return nil
}
