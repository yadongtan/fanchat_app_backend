package message

import (
	"fantastic_chat/server/database"
	"fmt"
	"time"
)

// SignInMessage 登录
type SignInMessage struct {
	TTid        int    `json:"ttid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Ip          string `json:"-"`
	DeviceModel string `json:"deviceModel"`
	DeviceName  string `json:"deviceName"`
	DeviceType  string `json:"deviceType"`
}

func (this *SignInMessage) Invoke() Message {
	userAccount := &database.UserAccount{}
	database.GetDB().Where("ttid = ?", this.TTid).First(userAccount)
	fmt.Printf("用户登录请求: 用户名{%s}, 密码{%s} \n", this.Username, this.Password)
	if this.TTid == userAccount.TTid && this.Password == userAccount.Password {

		// 添加登录日志
		ipDetails := database.GetIpDetails(this.Ip)

		signinLog := &database.UserSigninLog{
			TTid:        this.TTid,
			Type:        "Online",
			Ctime:       time.Now().Format("2006-01-02 15:04:05"),
			Ip:          this.Ip,
			Province:    ipDetails.Province,
			City:        ipDetails.City,
			Region:      ipDetails.Region,
			Addr:        ipDetails.Addr,
			DeviceModel: this.DeviceModel,
			DeviceName:  this.DeviceName,
			DeviceType:  this.DeviceType,
		}

		database.GetDB().Create(signinLog)
		return AckMessageOk("登录成功", userAccount.Username)
	} else {
		return AckMessageFailed("账号或密码错误", nil)
	}
}
