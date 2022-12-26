package message

import "fantastic_chat/server/database"

//添加朋友的消息
type AddFriendMessage struct {
	TTid       int `json:"ttid" gorm:"Column:ttid"`               //我的ttid
	FriendTTid int `json:"friend_ttid" gorm:"Column:friend_ttid"` //朋友的ttid
}

func (this *AddFriendMessage) Invoke() Message {
	// 向数据库中添加这条记录
	db := database.GetDB().Table("user_friend").Create(this)
	fri := &AddFriendMessage{
		this.FriendTTid,
		this.TTid,
	}
	db2 := database.GetDB().Table("user_friend").Create(fri)
	if db.Error != nil || db2.Error != nil {
		return AckMessageFailed("添加好友失败", nil)
	}
	return AckMessageOk("添加好友成功", nil)
}
