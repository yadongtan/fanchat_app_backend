package message

import (
	"encoding/json"
	"fantastic_chat/server/database"
)

//添加朋友的消息
type GetAIFriendsListMessage struct {
	TTid int `json:"ttid"`
}

func (this *GetAIFriendsListMessage) Invoke() Message {
	var list []database.OpenAIAccount
	r := database.GetDB().Table("openai_account as t1").Select(
		"t1.ai_ttid as ai_ttid ,t1.model as model, t1.name as name, t1.content as content, t1.ai_type as ai_type, t1.ctime as ctime").
		Joins("inner join user_friend as t2 on t1.ai_ttid = t2.friend_ttid").
		Where("t2.ttid = ? ", this.TTid).
		Find(&list)
	if r.Error != nil {
		return AckMessageFailed("获取AI列表失败", nil)
	}
	j, _ := json.Marshal(list)
	return AckMessageOk("获取AI列表成功", string(j))
}
