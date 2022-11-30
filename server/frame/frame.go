package frame

import (
	"fantastic_chat/server/encrypt"
	"fantastic_chat/server/message"
	"fantastic_chat/server/serialize"
	"fantastic_chat/server/utils"
	"fmt"
	"sync"
	"sync/atomic"
)

var ServerVersion = 1
var frameIdIncrement uint32 = 1

var wg sync.WaitGroup

// 原子操作版加函数
func getFrameId() int {
	wg.Add(1)
	defer wg.Done()
	return int(atomic.AddUint32(&frameIdIncrement, 1))
}

type Frame struct {
	FrameLen      int         //帧长
	Version       int         //版本号
	FrameId       int         //帧id
	FrameType     int         //帧类型
	SerializeType int         //序列化类型
	EncryptType   int         //消息体加密类型
	Payload       interface{} //消息体
}

// 将加密后的帧转换为字节
func GenerateFrameByesWithFrame(frame *Frame) []byte {
	bytes := make([]byte, 4*6)
	utils.CastIntToBytes(bytes, 0, frame.FrameLen)
	utils.CastIntToBytes(bytes, 4, frame.Version)
	utils.CastIntToBytes(bytes, 8, frame.FrameId)
	utils.CastIntToBytes(bytes, 12, frame.FrameType)
	utils.CastIntToBytes(bytes, 16, frame.SerializeType)
	utils.CastIntToBytes(bytes, 20, frame.EncryptType)
	bytes = append(bytes, frame.Payload.([]byte)...)
	return bytes
}

//生成帧并加密转换为bytes, 一些参数用默认值
func GenerateFrameBytesDefault(frameId int, payload interface{}) []byte {
	return GenerateFrameBytes(frameId, message.GetMessageTypeByInterface(payload), payload, encrypt.AESEncryptType, serialize.JsonSerializeType)
}

//生成帧并加密并转换为bytes
func GenerateFrameBytes(frameId int, frameType int, payload interface{}, encryptType int, serializeType int) []byte {
	//序列化
	serializedPayload := serialize.Serialize(serializeType, payload)
	// 加密
	encryptPayload := encrypt.Encrypt(serializedPayload, encryptType)
	payloadBytes := []byte(encryptPayload)
	version := ServerVersion
	frameLen := 4*6 + len(payloadBytes)

	bytes := make([]byte, 4*6)
	utils.CastIntToBytes(bytes, 0, frameLen)
	utils.CastIntToBytes(bytes, 4, version)
	utils.CastIntToBytes(bytes, 8, frameId)
	utils.CastIntToBytes(bytes, 12, frameType)
	utils.CastIntToBytes(bytes, 16, serializeType)
	utils.CastIntToBytes(bytes, 20, encryptType)
	bytes = append(bytes, payloadBytes...)
	fmt.Printf("加密后的数据:%v\n", payloadBytes)
	return bytes
}

func GenerateMessageFrame(msg interface{}) *Frame {
	f := GenerateFrame(getFrameId(), message.GetMessageTypeByInterface(msg), msg, encrypt.AESEncryptType, serialize.JsonSerializeType)
	return f
}

func GenerateAckFrame(from *Frame, ackMessage interface{}) *Frame {
	return GenerateFrame(from.FrameId, message.AckFrameType, ackMessage, from.EncryptType, from.SerializeType)
}

//生成帧并加密
func GenerateFrame(frameId int, frameType int, payload interface{}, encryptType int, serializeType int) *Frame {
	//序列化
	serializedPayload := serialize.Serialize(serializeType, payload)
	// 加密
	encryptPayload := encrypt.Encrypt(serializedPayload, encryptType)
	payloadBytes := []byte(encryptPayload)
	version := ServerVersion
	frameLen := 4*6 + len(payloadBytes)

	f := &Frame{
		frameLen,
		version,
		frameId,
		frameType,
		serializeType,
		encryptType,
		payloadBytes,
	}
	return f
}

//解析帧并解密
func ResolveFrame(bytes []byte) *Frame {
	frameLen := utils.CastBytesToInt(bytes, 0)
	version := utils.CastBytesToInt(bytes, 4)
	frameId := utils.CastBytesToInt(bytes, 8)
	frameType := utils.CastBytesToInt(bytes, 12)
	serializeType := utils.CastBytesToInt(bytes, 16)
	encryptType := utils.CastBytesToInt(bytes, 20)
	payload := bytes[24:]
	// fmt.Printf("获取到的完整数据:%v\n", bytes)
	// fmt.Printf("获取到的加密后的数据:%v\n", Payload)
	//解密
	de := encrypt.Decrypt(string(payload), encryptType)
	// 反序列化
	msg := serialize.UnserializeByType(serializeType, []byte(de), frameType)

	frame := Frame{
		frameLen, version, frameId, frameType, serializeType, encryptType, msg,
	}

	fmt.Printf("解析后的消息体 : %s\n", frame)
	return &frame
}
