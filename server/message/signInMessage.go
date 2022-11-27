package message

import "fmt"

// SignInMessage 登录
type SignInMessage struct {
	Uid      int
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *SignInMessage) Invoke() Message {
	fmt.Printf("用户登录请求: 用户名{%s}, 密码{%s} \n", this.Username, this.Password)
	this.Uid = 2209931449
	if this.Username == "yadong" && this.Password == "yadong123456" {
		return AckMessageOk
	} else {
		return AckMessageFailed
	}
}
