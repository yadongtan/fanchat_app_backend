package database

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
)

type UserAccount struct {
	TTid     int `gorm:"Column:ttid;PRIMARY_KEY;AUTO_INCREMENT;Column:ttid"`
	Username string
	Password string
	Ctime    string `gorm:"Column:ctime"`
}

type UserSigninLog struct {
	LogId       int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:log_id"`
	TTid        int    `gorm:"Column:ttid"`
	Type        string //online/offline
	Ctime       string `gorm:"Column:ctime"` //2002-01-01 01:01:01
	Ip          string
	Province    string
	City        string
	Region      string
	Addr        string
	DeviceModel string //设备模型
	DeviceName  string // 设备名称
	DeviceType  string // 设备类型
}

type IpDetails struct {
	Ip       string `json:"ip"`
	Province string `json:"province"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Addr     string `json:"addr"`
}

func GetIpDetails(ip string) *IpDetails {
	ipDetails := &IpDetails{}
	resp, err := http.Get("https://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true")
	if err != nil {
		return ipDetails
	}
	defer resp.Body.Close()
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("GetIpDetails() 获取ip详细信息时出错, err: ", err)
			return ipDetails
		}
		err = json.Unmarshal([]byte(ConvertGBK2Str(string(body))), ipDetails)
		if err != nil {
			fmt.Println("GetIpDetails() 获取ip详细信息时出错, err: ", err)
			return ipDetails
		}
		return ipDetails
	}
	return ipDetails
}

func ConvertStr2GBK(str string) string {
	//将utf-8编码的字符串转换为GBK编码
	ret, err := simplifiedchinese.GBK.NewEncoder().String(str)
	if err == nil {
		return ret //如果转换失败返回空字符串
	}
	//如果是[]byte格式的字符串，可以使用Bytes方法
	b, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	return string(b)
}

func ConvertGBK2Str(gbkStr string) string {
	//将GBK编码的字符串转换为utf-8编码
	ret, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err == nil {
		return ret //如果转换失败返回空字符串
	}
	//如果是[]byte格式的字符串，可以使用Bytes方法
	b, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkStr))
	return string(b)
}
