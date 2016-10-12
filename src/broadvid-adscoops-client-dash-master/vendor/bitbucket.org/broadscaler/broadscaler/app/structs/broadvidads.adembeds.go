package structs

type AdEmbedSave struct {
	ID                 int64  `form:"id"`
	AdID               string `form:"ad_id" gorm:"column:ad_id"`
	Type               string `form:"adType" sql:"-"`
	EmbedID            string `form:"embedID" gorm:"column:embed_id"`
	Sort               string `form:"sort"`
	AdEmbedLabel       string `form:"adEmbedLabel"`
	CustomField1       string `form:"customField1" gorm:"column:custom_field_1"`
	CustomField2       string `form:"customField2" gorm:"column:custom_field_2"`
	CustomField3       string `form:"customField3" gorm:"column:custom_field_3"`
	CustomField4       string `form:"customField4" gorm:"column:custom_field_4"`
	CustomField5       string `form:"customField5" gorm:"column:custom_field_5"`
	CustomFieldLabel1  string `sql:"-"`
	CustomFieldLabel2  string `sql:"-"`
	CustomFieldLabel3  string `sql:"-"`
	CustomFieldLabel4  string `sql:"-"`
	CustomFieldLabel5  string `sql:"-"`
	ForceSkip          string `form:"forceSkip"`
	ForceSkipKey       string `form:"forceSkipKey"`
	ForceSkipValue     string `form:"forceSkipValue"`
	UrlsWhitelist      string `form:"urls_whitelist"`
	UrlsBlacklist      string `form:"urls_blacklist"`
	CountriesWhitelist string `form:"countries_whitelist"`
	Pause              bool
}

type AdEmbed struct {
	AdEmbedSave
	EmbedLabel string
	Type       string
}

func (ae AdEmbed) TableName() string {
	return "ads_embeds"
}

func (ae AdEmbedSave) TableName() string {
	return "ads_embeds"
}

type AdEmbeds []AdEmbedSave

func (ae *AdEmbeds) FindAll() error {
	return BroadvidDB.Table("ads_embeds").Find(&ae).Error
}

func (ae *AdEmbedSave) Find(id string) error {
	return BroadvidDB.Table("ads_embeds").Find(&ae, id).Error
}

func (ae *AdEmbedSave) Remove(uid uint) {
	BroadvidDB.Table("ads_embeds").Delete(&ae)
}

func (ae *AdEmbedSave) PauseToggle(uid uint) {
	BroadvidDB.Table("ads_embeds").Where("id = ?", ae.ID).Find(&ae)
	ae.Pause = !ae.Pause

	BroadvidDB.Table("ads_embeds").Save(&ae)
}

func (ae *AdEmbedSave) Copy(uid uint) {
	var aecopy AdEmbedSave
	BroadvidDB.Table("ads_embeds").Where("id = ?", ae.ID).Find(&aecopy)

	aecopy.ID = 0
	aecopy.AdEmbedLabel = aecopy.AdEmbedLabel + " Copy"
	BroadvidDB.Table("ads_embeds").Save(&aecopy)
}

func (ae *AdEmbedSave) Save() error {
	var isNew = false
	var newType = ""

	if ae.ID == 0 {
		isNew = true
		newType = ae.Type
	}

	err := BroadvidDB.Table("ads_embeds").Save(&ae).Error

	if err != nil {
		return err
	}

	if isNew {
		BroadvidDB.Exec("UPDATE ads_embeds SET type = ? WHERE id = ?", newType, ae.ID)
	}

	return nil
}
