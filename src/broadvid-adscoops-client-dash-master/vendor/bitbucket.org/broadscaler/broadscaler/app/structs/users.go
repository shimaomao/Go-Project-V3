package structs

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/emailer"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email          string `json:"Email"`
	Name           string `json:"Name"`
	Userid         string `json:"UserID" sql:"unique_index"`
	Avatarurl      string `json:"AvatarURL"`
	Location       string `json:"Location"`
	LoginAllowed   bool
	IsAdscoopsUser bool
}

func (u *User) Find(id string) error {
	return AdscoopsDB.Find(&u, id).Error
}

func (u *User) Save() error {
	return AdscoopsDB.Save(&u).Error
}

type Users []User

func (us *Users) FindAll() error {
	return AdscoopsDB.Table("users").Find(&us).Error
}

func (us Users) MessageAdscoopsUsers(title, message string) (err error) {
	var usrs []User
	if err = AdscoopsDB.Where("is_adscoops_user = 1").Find(&usrs).Error; err != nil {
		log.Errorf("Cannot find adscoop users: %s", err)
		return
	}

	if len(usrs) == 0 {
		log.Warnf("No users found, so not sending a message")
		return
	}

	var emails []string

	for _, u := range usrs {
		emails = append(emails, u.Email)
	}

	var email emailer.Emailer

	if err = email.Send(title, message, emails); err != nil {
		return err
	}

	return

}

func (u *User) SaveJSON(user []byte) {
	var uj *User
	json.Unmarshal(user, &uj)
	err := AdscoopsDB.Where("email = ?", uj.Email).Find(&u).Error
	if err == nil {
		AdscoopsDB.Table("users").Save(&uj)
		u = uj
	}
}

func (u *User) FindById(id uint) {
	AdscoopsDB.Where("login_allowed = 1").Find(&u, id)
}

type UserAdscoopsClientSetting struct {
	gorm.Model
	ClientID     uint
	UserID       uint
	CampaignSort string
	ShowInfo     bool
	ClientOrder  uint
}

func (u *UserAdscoopsClientSetting) Save() error {
	if err := AdscoopsDB.Unscoped().Delete(UserAdscoopsClientSetting{}, "user_id = ? AND client_id = ?", u.UserID, u.ClientID).Error; err != nil {
		log.Errorf("Error deleting settings: %s", err)
	}
	return AdscoopsDB.Save(&u).Error
}
