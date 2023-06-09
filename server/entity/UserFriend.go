package entity

import "fantastic_chat/server/database"

//添加朋友
type UserFriend struct {
	TTid       int `gorm:"Column:ttid"`        //我的ttid
	FriendTTid int `gorm:"Column:friend_ttid"` //朋友的ttid
}

func AddFriend(ttid int, friendTTid int) error {
	my := UserFriend{TTid: ttid, FriendTTid: friendTTid}
	fri := UserFriend{TTid: friendTTid, FriendTTid: ttid}
	// 向数据库中添加这条记录
	db := database.GetDB().Table("user_friend").Create(my)
	db2 := database.GetDB().Table("user_friend").Create(fri)
	if db.Error != nil || db2.Error != nil {
		database.GetDB().Table("user_friend").Delete(my)
		database.GetDB().Table("user_friend").Delete(fri)
		return db.Error
	}
	return nil
}
