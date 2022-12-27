package message

import (
	"encoding/json"
	"fmt"
)

// JsonSerializeType json序列化方式
var JsonSerializeType = 1

func Unserialize(serializeType int, payload []byte, v interface{}) interface{} {
	switch serializeType {
	case JsonSerializeType:
		msg := json.Unmarshal(payload, v)
		return msg
	}
	return nil
}

//反序列化
func UnserializeByType(serializeType int, payload []byte, messageType int) interface{} {
	fmt.Printf("反序列化前:%v \n", string(payload))
	switch serializeType {
	case JsonSerializeType:
		msg := GetMessageByType(messageType)
		err := json.Unmarshal(payload, msg)
		if err != nil {
			fmt.Printf("反序列化失败 ! err : %v\n", err)
		}
		fmt.Printf("反序列化后:%v \n", msg)
		return msg
	}

	return nil
}

func Serialize(serializeType int, data interface{}) string {
	switch serializeType {
	case JsonSerializeType:
		payload, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Json Serialize failed! err: %v\n", err)
		}
		return string(payload)
	}
	return ""
}
