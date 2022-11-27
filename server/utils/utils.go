package utils

func CastIntToBytes(bytes []byte, startIndex int, data int) {
	bytes[startIndex] = byte(data >> 24)
	bytes[startIndex+1] = byte(data >> 16)
	bytes[startIndex+2] = byte(data >> 8)
	bytes[startIndex+3] = byte(data)
}

func CastBytesToInt(bytes []byte, startIndex int) int {
	data := 0

	data += int(bytes[startIndex]) << 24
	data += int(bytes[startIndex+1]) << 16
	data += int(bytes[startIndex+2]) << 8
	data += int(bytes[startIndex+3])

	return data
}
