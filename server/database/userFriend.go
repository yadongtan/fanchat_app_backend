package database

//添加朋友的消息
type UserFriend struct {
	TTid       int `json:"ttid" gorm:"Column:ttid"`               //我的ttid
	FriendTTid int `json:"friend_ttid" gorm:"Column:friend_ttid"` //朋友的ttid
}
