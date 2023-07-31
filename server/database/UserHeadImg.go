package database

import "gorm.io/gorm"

type UserHeadImg struct {
	TTid  int    `json:"ttid" gorm:"Column:ttid;PRIMARY_KEY;"`
	ImgId string `json:"ImgId" gorm:"Column: img_id"`
}

func (this *UserHeadImg) UpdateOrInsert() error {
	old := UserHeadImg{}
	result := GetDB().Model(this).Where("ttid = ?", this.TTid).First(&old)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			result := GetDB().Model(this).Create(this) //没有则创建
			if result.Error != nil {
				return result.Error
			}
		}
	} else {
		// 找到了, 有则更新
		result := GetDB().Model(this).Updates(this)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
