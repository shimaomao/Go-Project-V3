package structs

import "strings"

type VideoThemes []VideoTheme

func (v *VideoThemes) FindAll() error {
	return BroadvidDB.Table("video_themes").Find(&v).Error
}

type VideoTheme struct {
	ID   int64  `form:"id"`
	Code string `form:"code"`
	Name string `form:"name"`
}

func (v *VideoTheme) Find(id string) error {
	return BroadvidDB.Find(&v, id).Error
}

func (v *VideoTheme) Save() error {
	if !strings.Contains(v.Code, "send.adscoops.com/tracking.js") {
		v.Code = strings.Replace(v.Code, "</body>", `<script type="text/javascript" src="http://send.adscoops.com/tracking.js"></script></body>`, 1)
	}
	return BroadvidDB.Save(&v).Error
}
