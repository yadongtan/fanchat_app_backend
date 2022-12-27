package message

import (
	"fantastic_chat/server/database"
	"time"
)

// SignUpMessage 注册
type SignUpMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func (this *SignUpMessage) Invoke() Message {

	userAccount := &database.UserAccount{}

	userAccount.Username = this.Username
	userAccount.Password = this.Password
	userAccount.Phone = this.Phone
	userAccount.Ctime = time.Now().Format("2006-01-02 15:04:05")

	// 创建账号S
	database.GetDB().Create(userAccount)

	// 查询刚刚创建账号的id
	database.GetDB().Where("username = ? and password = ? and ctime = ?",
		userAccount.Username, userAccount.Password, userAccount.Ctime).First(userAccount)

	return AckMessageOk("注册成功", userAccount.TTid)

}
