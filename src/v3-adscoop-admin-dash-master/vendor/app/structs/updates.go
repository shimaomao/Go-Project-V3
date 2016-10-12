package structs

import "github.com/jinzhu/gorm"

type Updates struct {
	gorm.Model
	UserID  uint
	Title   string
	Message string
	Product uint // 1: adscoops, 2: broadvid ads, 3: broadvid vids
}

func (u *Updates) Save() {
	AdscoopsDB.Save(&u)
}
