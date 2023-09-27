package message

import (
	"fantastic_chat/server/database"
	"os"
	"strings"
	"time"
)

type UpdateUserHeadImgMessage struct {
	TTid    int    `json:"ttid" `
	ImgName string `json:img_name`
	Img     []byte `json:"img"`
}

func (this *UpdateUserHeadImgMessage) Invoke() Message {
	if this.ImgName != "" {
		strs := strings.Split(this.ImgName, ".")
		if len(strs) > 1 {
			suffix := strs[len(strs)-1]
			this.ImgName = string(rune(this.TTid)) + time.Now().Format("20060102150405") + "." + suffix
		} else {
			//没有后缀名
			this.ImgName = string(rune(this.TTid)) + time.Now().Format("20060102150405")
		}
	} else {
		// 图片没有设置名字
		this.ImgName = string(rune(this.TTid)) + time.Now().Format("20060102150405")
	}
	// 图片名字有了, 保存图片到本地
	saveDir := "./img/head/"
	savePath := saveDir + this.ImgName
	file, err := os.Create(savePath)
	if err != nil {
		return AckMessageFailed("上传图片失败", nil)
	}
	for n := 0; n < len(this.Img); {
		cnt, _ := file.Write(this.Img)
		n += cnt
	}
	dbUserHeadImg := &database.UserHeadImg{
		TTid:  this.TTid,
		ImgId: this.ImgName,
	}
	err = dbUserHeadImg.UpdateOrInsert()
	if err != nil {
		return AckMessageFailed("上传图片失败", nil)
	}
	return AckMessageOk("上传图片成功", this.ImgName)

}
