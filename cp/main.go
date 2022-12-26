package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func main() {
	pwHash := "ea536bb7d8e51acc92f280fdea7798f21819dd13a9a8e7bc78"
	data, err := ioutil.ReadFile("./cp/password.txt")
	if err != nil {
		panic(err)
	}
	pw := string(data)
	str := strings.Split(pw, "\r\n")
	// 遍历每一个字符串
	for _, s := range str {
		fmt.Println("比较密码: ", s)
		if isPsCorrect(s, pwHash) {
			fmt.Println("正确密码: ", s)
			break
		}
	}
	isPsCorrect("kkl4Ib", pwHash)
}

func sha1WithSalt(pw string, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pw)
	return string(h.Sum(nil))
}

func isPsCorrect(pw string, hash string) bool {
	if len(hash) == 50 {
		saltedHash := hash[0:40]
		salt := hash[40:50]
		fmt.Println("...")
		temp := sha1WithSalt(pw, salt)
		fmt.Printf("%x\n", []byte(temp))
		if temp == saltedHash {
			return true
		} else {
			return false
		}
	}
	return false
}
