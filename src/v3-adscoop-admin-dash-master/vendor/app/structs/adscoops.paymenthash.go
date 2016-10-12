package structs

import "github.com/jinzhu/gorm"

type AdscoopPaymentHash struct {
	gorm.Model
	Hash     string
	ClientID uint
}
