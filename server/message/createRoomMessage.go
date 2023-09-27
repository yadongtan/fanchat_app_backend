package message

import "fantastic_chat/server/database"

type CreateRoomMessage struct {
	TTid     int    `json:"ttid"`     //我的ttid
	Username string `json:"username"` //用户名
	Model    string `json:"model"`    //模型名称
}

func (this *CreateRoomMessage) Invoke() Message {
	rMap := GetRoomMap()
	ua := database.UserAccount{TTid: this.TTid, Username: this.Username}
	rMap.AddUserToRoom(this.TTid, ua)

	return AckMessageOk("创建房间成功", nil)
}
