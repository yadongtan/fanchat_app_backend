package message

import (
	"fantastic_chat/server/database"
)

//开启一个与OpenAIChat模型的聊天
type CreateChatMessage struct {
	TTid     int    `json:"ttid"`     //我的ttid
	Username string `json:"username"` //用户名
	Model    string `json:"model"`    //模型名称
}

func (this *CreateChatMessage) Invoke() Message {
	// 为用户生成一个AI的账号, 虚拟出来
	err := database.CreateOpenAIChatAccount(this.TTid, "text", this.Model)
	if err != nil {
		return AckMessageFailed("生成AI账号失败!", err.Error())
	}
	return AckMessageOk("添加好友成功", nil)
}
