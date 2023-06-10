package message

import (
	"fantastic_chat/server/encrypt"
	"fantastic_chat/server/utils"
	"fmt"
)

var ServerVersion = 1

type Frame struct {
	FrameLen      int         //帧长
	Version       int         //版本号
	FrameType     int         //帧类型
	SerializeType int         //序列化类型
	EncryptType   int         //消息体加密类型
	FrameId       string      //帧id
	Payload       interface{} //消息体
}

// 将加密后的帧转换为字节
func GenerateFrameByesWithFrame(frame *Frame) []byte {

	return CastFrameToByte(frame.FrameLen, frame.Version, frame.FrameType, frame.SerializeType, frame.EncryptType, frame.FrameId, frame.Payload)

}

//生成帧并加密转换为bytes, 一些参数用默认值
func GenerateFrameBytesDefault(frameId string, payload interface{}) []byte {
	return GenerateFrameBytes(frameId, GetMessageTypeByInterface(payload), payload, encrypt.AESEncryptType, JsonSerializeType)
}

func CastFrameToByte(frameLen int, version int, frameType int, serializeType int, encryptType int, frameId string, payload interface{}) []byte {
	bytes := make([]byte, 5*4)
	utils.CastIntToBytes(bytes, 0, frameLen)
	utils.CastIntToBytes(bytes, 4, version)
	utils.CastIntToBytes(bytes, 8, frameType)
	utils.CastIntToBytes(bytes, 12, serializeType)
	utils.CastIntToBytes(bytes, 16, encryptType)
	// frameId
	bytes = append(bytes, ([]byte)(frameId)...)
	// payload
	bytes = append(bytes, payload.([]byte)...)
	return bytes
}

//生成帧并加密并转换为bytes
func GenerateFrameBytes(frameId string, frameType int, payload interface{}, encryptType int, serializeType int) []byte {
	//序列化
	serializedPayload := Serialize(serializeType, payload)
	// 加密
	encryptPayload := encrypt.Encrypt(serializedPayload, encryptType)
	payloadBytes := []byte(encryptPayload)
	version := ServerVersion
	frameLen := 4*6 + len(payloadBytes)

	bytes := CastFrameToByte(frameLen, version, frameType, serializeType, encryptType, frameId, payload)

	fmt.Printf("加密后的数据:%v\n", payloadBytes)
	return bytes
}

func GenerateMessageFrame(frameId string, msg interface{}) *Frame {
	f := GenerateFrame(frameId, GetMessageTypeByInterface(msg), msg, encrypt.AESEncryptType, JsonSerializeType)
	return f
}

func GenerateAckFrame(from *Frame, ackMessage interface{}) *Frame {
	return GenerateFrame(from.FrameId, AckFrameType, ackMessage, from.EncryptType, from.SerializeType)
}

//生成帧并加密
func GenerateFrame(frameId string, frameType int, payload interface{}, encryptType int, serializeType int) *Frame {
	//序列化
	serializedPayload := Serialize(serializeType, payload)
	// 加密
	encryptPayload := encrypt.Encrypt(serializedPayload, encryptType)
	payloadBytes := []byte(encryptPayload)
	version := ServerVersion
	frameLen := 4*5 + len(payloadBytes) + len(frameId)

	f := &Frame{
		frameLen,
		version,
		frameType,
		serializeType,
		encryptType,
		frameId,
		payloadBytes,
	}
	return f
}

//解析帧并解密
func ResolveFrame(bytes []byte) *Frame {
	frameLen := utils.CastBytesToInt(bytes, 0)
	version := utils.CastBytesToInt(bytes, 4)
	frameType := utils.CastBytesToInt(bytes, 8)
	serializeType := utils.CastBytesToInt(bytes, 12)
	encryptType := utils.CastBytesToInt(bytes, 16)

	frameId := string(bytes[20:39])
	payload := bytes[39:]
	// fmt.Printf("获取到的完整数据:%v\n", bytes)
	// fmt.Printf("获取到的加密后的数据:%v\n", Payload)
	//解密
	de := encrypt.Decrypt(string(payload), encryptType)
	// 反序列化
	msg := UnserializeByType(serializeType, []byte(de), frameType)

	frame := Frame{
		frameLen, version, frameType, serializeType, encryptType, frameId, msg,
	}

	fmt.Printf("解析后的消息体 : %s\n", frame)
	return &frame
}
