package structs

type AdUrlGroups []AdUrlGroup

func (a *AdUrlGroups) FindAll() error {
	return BroadvidDB.Table("ad_url_groups").Find(&a).Error
}

type AdUrlGroup struct {
	ID   int64    `form:"id"`
	Name string   `form:"name"`
	Urls []string `form:"url[]" sql:"-"`
}

func (a *AdUrlGroup) Find(id string) error {
	err := BroadvidDB.Table("ad_url_groups").Find(&a, id).Error

	if err != nil {
		return err
	}

	var aus []AdUrl

	err = BroadvidDB.Where("ad_url_group_id = ?", id).Find(&aus).Error

	if err != nil {
		return err
	}

	for _, u := range aus {
		a.Urls = append(a.Urls, u.Url)
	}

	return nil
}

func (a *AdUrlGroup) Save() error {
	err := BroadvidDB.Table("ad_url_groups").Save(&a).Error

	if err != nil {
		return err
	}

	var au AdUrl

	err = BroadvidDB.Where("ad_url_group_id = ?", a.ID).Delete(&au).Error

	if err != nil {
		return err
	}

	for _, u := range a.Urls {
		var au AdUrl
		au.Url = u
		au.AdUrlGroupID = a.ID
		BroadvidDB.Save(&au)
	}

	return err
}

type AdUrl struct {
	ID           int64  `form:"id"`
	Url          string `form:"url"`
	AdUrlGroupID int64  `form:"ad_url_group_id"`
}
