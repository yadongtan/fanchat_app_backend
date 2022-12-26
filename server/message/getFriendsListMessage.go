package message

import (
	"encoding/json"
	"fantastic_chat/server/channel"
	"fantastic_chat/server/database"
)

//添加朋友的消息
type GetFriendsListMessage struct {
	TTid int `json:"ttid" gorm:"Column:ttid"` //我的ttid
}

func (this *GetFriendsListMessage) Invoke() Message {
	// 向数据库中添加这条记录
	var list []database.UserFriendDetails
	// select t1.ttid,t1.friend_ttid , t2.username from user_friend as t1 inner join user_accounts as t2 on t1.friend_ttid = t2.ttid;
	r := database.GetDB().Table("user_friend as t1").Select("t1.ttid as ttid ,t1.friend_ttid as friend_ttid  , t2.username as friend_username").
		Joins("inner join user_accounts as t2 on t1.friend_ttid = t2.ttid").
		Where("t1.ttid = ?", this.TTid).
		Find(&list)
	if r.Error != nil {
		return AckMessageFailed("查询好友列表失败", nil)
	}
	for i := 0; i < len(list); i++ {
		if channel.Cs.OnlineUserMap[list[i].TTid] != nil {
			list[i].Status = database.OnlineStatus
		} else {
			list[i].Status = database.OfflineStatus
		}
	}
	j, _ := json.Marshal(list)
	return AckMessageOk("查询好友列表成功", string(j))
}
