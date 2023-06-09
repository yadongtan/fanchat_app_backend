package message

import (
	"fantastic_chat/server/database"
	"fmt"
)

//添加朋友的消息
type FriendDirectMessage struct {
	FromTTid int    `json:"from_ttid" gorm:"Column:from_ttid"` //我的ttid
	DestTTid int    `json:"dest_ttid" gorm:"Column:dest_ttid"` //朋友的ttid
	Time     int64  `json:"time" gorm:"Column:time"`           //消息发出时间
	Text     string `json:"text" gorm:"Column:text"`           //消息内容
}

func (this *FriendDirectMessage) Invoke() Message {
	// 向数据库中添加这条记录
	db := database.GetDB().Table("direct_msg").Create(this)
	// 查看好友是否在线, 如果在, 转发该消息给该好友

	c, ok := OnlineUserChannelMap.Load(this.DestTTid)
	if ok {
		go c.(*Channel).Write(this)
	}

	if db.Error != nil {
		return AckMessageFailed("发送消息失败", nil)
	}
	return AckMessageOk("发送消息成功", nil)
}

// 获取指定时间之后的私发消息
func GetDirectMsg(fromTime float32, ttid int) ([]FriendDirectMessage, error) {
	var list []FriendDirectMessage

	db := database.GetDB().Table("direct_msg").Where("time > ? and (from_ttid = ? or dest_ttid = ? )", fromTime, ttid, ttid).Find(&list)

	if len(list) == 0 {
		fmt.Printf("查询离线期间好友消息为空\n")
	} else {
		for _, msg := range list {
			fmt.Printf("查询离线期间好友消息: %v\n", msg)
		}
	}

	if db.Error != nil {
		return nil, db.Error
	}
	return list, nil
}
