package database

//添加朋友的消息
type UserFriendDetails struct {
	TTid           int    `json:"ttid" gorm:"Column:ttid"`                 //我的ttid
	FriendTTid     int    `json:"friend_ttid" gorm:"Column:friend_ttid"`   //朋友的ttid
	FriendUsername string `json:"friend_username", gorm:"friend_username"` //朋友的username
	Status         string `json:"status"`                                  //朋友的在线状态
}

var OnlineStatus = "[在线]"
var OfflineStatus = "[离线]"
