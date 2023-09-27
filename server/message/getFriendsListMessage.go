package message

import (
	"encoding/json"
	"fantastic_chat/server/database"
)

//添加朋友的消息
type GetFriendsListMessage struct {
	TTid int `json:"ttid" gorm:"Column:ttid"` //我的ttid
}

func (this *GetFriendsListMessage) Invoke() Message {
	var list []database.UserFriendDetails
	r := database.GetDB().Table("user_friend as t1").Select("t1.ttid as ttid ,t1.friend_ttid as friend_ttid  , t2.username as friend_username").
		Joins("inner join user_accounts as t2 on t1.friend_ttid = t2.ttid").
		Where("t1.ttid = ? ", this.TTid).
		Find(&list)

	if r.Error != nil {
		return AckMessageFailed("查询好友列表失败", nil)
	}

	for i := 0; i < len(list); i++ {
		_, ok := OnlineUserChannelMap.Load(list[i].FriendTTid)
		if ok {
			list[i].Status = database.OnlineStatus
		} else {
			list[i].Status = database.OfflineStatus
		}
	}
	j, _ := json.Marshal(list)
	return AckMessageOk("查询好友列表成功", string(j))
}
