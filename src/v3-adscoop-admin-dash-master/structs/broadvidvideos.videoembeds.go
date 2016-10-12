package structs

type VideoEmbeds []VideoEmbed

func (v *VideoEmbeds) FindAll() error {
	return BroadvidDB.Table("video_embeds").Find(&v).Error
}

type VideoEmbed struct {
	ID                int64  `form:"id"`
	Label             string `form:"label"`
	Embed             string `form:"embed"`
	CustomFieldLabel1 string `form:"customFieldLabel1" gorm:"column:custom_field_label_1"`
	CustomFieldLabel2 string `form:"customFieldLabel2" gorm:"column:custom_field_label_2"`
	CustomFieldLabel3 string `form:"customFieldLabel3" gorm:"column:custom_field_label_3"`
	CustomFieldLabel4 string `form:"customFieldLabel4" gorm:"column:custom_field_label_4"`
	CustomFieldLabel5 string `form:"customFieldLabel5" gorm:"column:custom_field_label_5"`
}

func (v *VideoEmbed) Find(id string) error {
	return BroadvidDB.Find(&v, id).Error
}

func (v *VideoEmbed) Save() error {
	return BroadvidDB.Save(&v).Error
}

func (v VideoEmbed) TableName() string {
	return "video_embeds"
}
