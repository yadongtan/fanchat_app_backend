package database

import (
	"errors"
	"time"
)

type OpenAIAccount struct {
	TTid    int    `json:"ttid" gorm:"Column:ttid"`
	AiTTid  int    `json:"ai_ttid" gorm:"Column:ai_ttid"`
	Model   string `json:"model" gorm:"Column:model"`
	Name    string `json:"name" gorm:"Column:name"`
	Content string `json:"content" gorm:"Column:content"`
	AIType  string `json:"ai_type" gorm:"Column:ai_type"`
	Ctime   string `json:"ctime" gorm:"Column:ctime"`
}

const (
	TextType  string = "text"
	ImageType string = "image"
)

const (
	GPT_3_5_Turbo string = "gpt-3.5-turbo"
)

func CreateOpenAIChatAccount(ttid int, aiType string, model string) error {

	aiAccount := OpenAIAccount{}
	aiAccount.TTid = ttid
	aiAccount.AIType = aiType
	aiAccount.Model = model

	aiAccount.Ctime = time.Now().Format("2006-01-02 15:04:05")

	aiAccount.Name = model + "_" + aiType + "_" + aiAccount.Ctime

	// 创建账号
	tx := GetDB().Table("openai_account").Create(&aiAccount)
	if tx.Error != nil {
		return tx.Error
	}

	GetDB().Table("openai_account").Where("ttid = ? and name = ?", aiAccount.TTid, aiAccount.Name).First(&aiAccount)
	if aiAccount.AiTTid == 0 {
		return errors.New("创建AI账号失败")
	}
	// 绑定与创建者的关系
	err := AddFriend(aiAccount.TTid, aiAccount.AiTTid)
	if err != nil {
		// 删除创建的账号
		GetDB().Table("openai_account").Delete(aiAccount)
		return err
	}
	return nil

}
