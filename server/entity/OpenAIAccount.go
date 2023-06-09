package entity

import (
	"errors"
	"fantastic_chat/server/database"
	"time"
)

type OpenAIAccount struct {
	TTid    int    `gorm:"Column:ttid""`
	AiTTid  int    `gorm:"Column:ai_ttid"`
	Model   string `gorm:"Column:model"`
	Name    string `gorm:"Column:name"`
	Content string `gorm:"Column:content"`
	Type    string `gorm:"Column:type"`
	Ctime   string `gorm:"Column:ctime"`
}

func CreateOpenAIChatAccount(ttid int, aiType string, model string) error {

	aiAccount := OpenAIAccount{}
	aiAccount.TTid = ttid
	aiAccount.Type = aiType
	aiAccount.Model = model

	aiAccount.Ctime = time.Now().Format("2006-01-02 15:04:05")

	aiAccount.Name = aiType + " " + model + " " + aiAccount.Ctime

	// 创建账号
	tx := database.GetDB().Table("openai_account").Create(&aiAccount)
	if tx.Error != nil {
		return tx.Error
	}

	database.GetDB().Table("openai_account").Where("ttid = ? and name = ?", aiAccount.TTid, aiAccount.Name).First(&aiAccount)
	if aiAccount.AiTTid == 0 {
		return errors.New("创建AI账号失败")
	}
	// 绑定与创建者的关系
	err := AddFriend(aiAccount.TTid, aiAccount.AiTTid)
	if err != nil {
		// 删除创建的账号
		database.GetDB().Table("openai_account").Delete(aiAccount)
		return err
	}
	return nil

}
