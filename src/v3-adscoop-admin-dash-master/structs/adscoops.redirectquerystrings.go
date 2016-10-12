package structs

import "github.com/jinzhu/gorm"

type RedirectQuerystring struct {
	gorm.Model
	RedirectID     uint
	QueryStringKey string
}

func (a RedirectQuerystring) TableName() string {
	return "adscoop_redirect_querystrings"
}

func (r *RedirectQuerystring) Save() error {
	return AdscoopsDB.Save(&r).Error
}
