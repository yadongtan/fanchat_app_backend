package message

import (
	"encoding/json"
	"fmt"
)

type GetRoomListMessage struct {
}

func (this *GetRoomListMessage) Invoke() Message {
	j, err := json.Marshal(GetRoomMap())
	if err != nil {
		fmt.Println(err)
		return AckMessageFailed("获取房间列表失败", nil)
	}
	return AckMessageOk("获取房间列表成功", string(j))
}
