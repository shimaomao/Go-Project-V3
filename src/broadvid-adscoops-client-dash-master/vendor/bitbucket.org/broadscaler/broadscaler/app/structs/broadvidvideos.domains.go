package structs

type VideoDomains []VideoDomain

func (v *VideoDomains) FindAll() error {
	return BroadvidDB.Table("video_domains").Find(&v).Error
}

type VideoDomain struct {
	ID              int64  `form:"id"`
	Host            string `form:"host"`
	GaID            string `form:"ga_id"`
	ThemeID         int64  `form:"theme_id"`
	DefaultRedirect int64  `form:"default_redirect"`
}

func (v *VideoDomain) Find(id string) error {
	return BroadvidDB.Find(&v, id).Error
}

func (v *VideoDomain) Save() error {
	return BroadvidDB.Save(&v).Error
}
