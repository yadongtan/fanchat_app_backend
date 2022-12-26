package database

import "gorm.io/gorm"

type UserDetails struct {
	TTid          string `json:"ttid" gorm:"Column:ttid;PRIMARY_KEY;"`
	Sex           string `json:"sex" gorm:"Column:sex"`
	Age           int    `json:"age" gorm:"Column:age"`
	Birthday      string `json:"birthday" gorm:"Column:birthday"`
	Location      string `json:"location" gorm:"Column:location"`
	Hometown      string `json:"hometown" gorm:"Column:hometown"`
	Constellation string `json:"constellation" gorm:"Column:constellation"`
}

func InsertUserDetails(details *UserDetails) (tx *gorm.DB) {
	return GetDB().Create(details)
}

func UpdateUserDetails(details *UserDetails) (tx *gorm.DB) {
	return GetDB().Model(details).Updates(details)
}
