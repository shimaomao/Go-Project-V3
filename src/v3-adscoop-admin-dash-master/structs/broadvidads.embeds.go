package structs

type Embeds []Embed

func (e *Embeds) FindAll() error {
	return BroadvidDB.Table("embeds").Find(&e).Error
}

func (a *Embed) GetID() int64 {
	return a.ID
}

func (e *Embed) Find(id string) error {
	return BroadvidDB.Table("embeds").Find(&e, id).Error
}

func (e *Embed) Save() error {
	return BroadvidDB.Table("embeds").Save(&e).Error
}

type Embed struct {
	ID                int64  `form:"id"`
	Label             string `form:"label" gorm:"column:embed_label"`
	Code              string `form:"code"`
	CodeExternal      string `form:"codeExternal"`
	Type              string `form:"type"`
	CustomField1Label string `form:"customField1Label" gorm:"column:custom_field_1_label"`
	CustomField2Label string `form:"customField2Label" gorm:"column:custom_field_2_label"`
	CustomField3Label string `form:"customField3Label" gorm:"column:custom_field_3_label"`
	CustomField4Label string `form:"customField4Label" gorm:"column:custom_field_4_label"`
	CustomField5Label string `form:"customField5Label" gorm:"column:custom_field_5_label"`
	Selected          bool   `sql:"-"`
}
