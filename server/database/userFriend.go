package database

//添加朋友的消息
type UserFriend struct {
	TTid       int `json:"ttid" gorm:"Column:ttid"`               //我的ttid
	FriendTTid int `json:"friend_ttid" gorm:"Column:friend_ttid"` //朋友的ttid
}

func AddFriend(ttid int, friendTTid int) error {
	my := UserFriend{TTid: ttid, FriendTTid: friendTTid}
	fri := UserFriend{TTid: friendTTid, FriendTTid: ttid}
	// 向数据库中添加这条记录
	db := GetDB().Table("user_friend").Create(my)
	db2 := GetDB().Table("user_friend").Create(fri)
	if db.Error != nil || db2.Error != nil {
		GetDB().Table("user_friend").Delete(my)
		GetDB().Table("user_friend").Delete(fri)
		return db.Error
	}
	return nil
}
