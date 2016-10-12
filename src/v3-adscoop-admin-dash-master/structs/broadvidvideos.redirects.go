package structs

type VideoRedirectors []VideoRedirector

func (v *VideoRedirectors) FindAll() error {
	return BroadvidDB.Table("video_redirector").Find(&v).Error
}

type VideoRedirector struct {
	ID             int64  `form:"id"`
	RssID          int64  `form:"rssId"`
	Key            string `form:"key"`
	StripParams    bool
	ClientRedirect bool
}

func (v VideoRedirector) TableName() string {
	return "video_redirector"
}

func (v *VideoRedirector) Find(id string) error {
	return BroadvidDB.Find(&v, id).Error
}

func (v *VideoRedirector) Save() error {
	return BroadvidDB.Save(&v).Error
}
