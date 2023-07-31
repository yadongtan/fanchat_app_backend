package encrypt

import "fmt"

var AESEncryptType = 1

func Encrypt(data string, encryptType int) []byte {
	switch encryptType {
	case AESEncryptType:
		en, err := AesEcpt.AesBase64Encrypt(data)
		if err != nil {
			fmt.Printf("encrypt error:%v\n", err)
			return nil
		}
		return []byte(en)
	}
	return nil
}

func Decrypt(data string, encryptType int) string {
	switch encryptType {
	case AESEncryptType:
		de, err := AesEcpt.AesBase64Decrypt(data)
		if err != nil {
			fmt.Printf("encrypt error:%v\n", err)
			return ""
		}
		return string(de)
	}
	return ""
}
