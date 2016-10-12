package structs

import "time"

type VideoRsss []VideoRss

func (v *VideoRsss) FindAll() error {
	return BroadvidDB.Table("video_rss").Find(&v).Error
}

type VideoRss struct {
	ID                int64
	RssUrl            string
	Type              string
	CustomField1      string `form:"customField1" gorm:"column:custom_field_1"`
	CustomField2      string `form:"customField2" gorm:"column:custom_field_2"`
	CustomField3      string `form:"customField3" gorm:"column:custom_field_3"`
	CustomField4      string `form:"customField4" gorm:"column:custom_field_4"`
	CustomField5      string `form:"customField5" gorm:"column:custom_field_5"`
	EmbedID           int64  `form:"embedId"`
	CreatedOn         time.Time
	LastFetched       time.Time
	Width             int64       `form:"width"`
	Height            int64       `form:"height"`
	Mute              int64       `form:"mute"`
	AutoPlay          bool        `form:"auto_play"`
	Label             string      `form:"label"`
	OneTimePull       bool        `form:"one_time_pull"`
	Items             []VideoItem `sql:"-"`
	OverrideThemeID   int64       `form:"override_theme_id"`
	AstMute           bool        `form:"ast_mute"`
	AstSeek           bool        `form:"ast_seek"`
	VerifyFireTimeout int64       `form:"verify_fire_timeout"`
	TodaysCount       int64       `form:"todays_count"`
}

func (v *VideoRss) Find(id string) error {
	return BroadvidDB.Table("video_rss").Find(&v, id).Error
}

func (v *VideoRss) Save() error {

	v.CustomField1 = cleanUpMacros(v.CustomField1)
	v.CustomField2 = cleanUpMacros(v.CustomField2)
	v.CustomField3 = cleanUpMacros(v.CustomField3)
	v.CustomField4 = cleanUpMacros(v.CustomField4)
	v.CustomField5 = cleanUpMacros(v.CustomField5)
	return BroadvidDB.Save(&v).Error
}

type VideoRssOverride struct {
	VideoRssID int64  `form:"video_rss_id"`
	ID         int64  `form:"id"`
	Mute       int64  `form:"mute"`
	AutoPlay   int64  `form:"auto_play"`
	Value      string `form:"value"`
}

func (v VideoRss) TableName() string {
	return "video_rss"
}
