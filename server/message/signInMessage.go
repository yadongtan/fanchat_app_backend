package message

import "fmt"

// SignInMessage 登录
type SignInMessage struct {
	TTid     int    `json:"ttid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *SignInMessage) Invoke() Message {
	fmt.Printf("用户登录请求: 用户名{%s}, 密码{%s} \n", this.Username, this.Password)

	if this.TTid == 123456 && this.Password == "123456" {
		return AckMessageOk(nil)
	} else if this.TTid == 220993 && this.Password == "220993" {
		return AckMessageOk(nil)
	} else if this.TTid == 859888 && this.Password == "859888" {
		return AckMessageOk(nil)
	} else {
		return AckMessageFailed(nil)
	}
}
