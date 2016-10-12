package app

import (
	"log"
	"time"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
)

type DefaultRetData struct {
	User *UserWithPolicy
}

type tomlConfig struct {
	SqlConnection       string
	PaymentHost         string
	StripeSecretKey     string
	StripePublishKey    string
	MailGunDomain       string
	MailGunApiKey       string
	MailGunPublicApikey string
	RedirHost           string
}

type UserWithPolicy struct {
	structs.AdscoopClientUser
	Policy structs.AdscoopUserPolicy
}

func (u UserWithPolicy) TableName() string {
	return "adscoop_client_users"
}

/*
	WeekdayStartHour                              uint   `form:"weekday_start_hour"`
	WeekdayEndHour                                uint   `form:"weekday_end_hour"`
	WeekendStartHour                              uint   `form:"weekend_start_hour"`
	WeekendEndHour                                uint   `form:"weekend_end_hour"`
*/

func (u *UserWithPolicy) IsBusinessHours() bool {
	if !u.Policy.EnableBusinessHours {
		return true
	}
	location, _ := time.LoadLocation("America/Los_Angeles")

	loc := time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(), 0, 0, location)

	log.Println("now", loc.Hour())

	if loc.Weekday() == 1 && (u.Policy.MondayBlackout || (u.Policy.MondayStartHour > loc.Hour() || u.Policy.MondayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Monday")
		return false
	}

	if loc.Weekday() == 2 && (u.Policy.TuesdayBlackout || (u.Policy.TuesdayStartHour > loc.Hour() || u.Policy.TuesdayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Tuesday")
		return false
	}

	if loc.Weekday() == 3 && (u.Policy.WednesdayBlackout || (u.Policy.WednesdayStartHour > loc.Hour() || u.Policy.WednesdayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Wednesday")
		return false
	}

	if loc.Weekday() == 4 && (u.Policy.ThursdayBlackout || (u.Policy.ThursdayStartHour > loc.Hour() || u.Policy.ThursdayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Thursday")
		return false
	}

	if loc.Weekday() == 5 && (u.Policy.FridayBlackout || (u.Policy.FridayStartHour > loc.Hour() || u.Policy.FridayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Friday")
		return false
	}

	if loc.Weekday() == 6 && (u.Policy.SaturdayBlackout || (u.Policy.SaturdayStartHour > loc.Hour() || u.Policy.SaturdayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Saturday")
		return false
	}

	if loc.Weekday() == 7 && (u.Policy.SundayBlackout || (u.Policy.SundayStartHour > loc.Hour() || u.Policy.SundayEndHour <= loc.Hour())) {
		log.Println("Out of business hours on Sunday")
		return false
	}

	return true
}
