package message

import "fantastic_chat/server/database"

type UpdateUserAccountMessage struct {
	UserAccount *database.UserAccount `json:"UserAccount"`
}

func (this *UpdateUserAccountMessage) Invoke() Message {
	if this.UserAccount == nil {
		return AckMessageFailed("failed", "更新失败")
	}
	database.GetDB().Model(this.UserAccount).Where("ttid = ?", this.UserAccount.TTid).Updates(this.UserAccount)
	return AckMessageOk("ok", nil)
}
