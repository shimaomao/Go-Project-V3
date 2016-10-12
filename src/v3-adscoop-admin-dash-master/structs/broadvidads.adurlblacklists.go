package structs

type AdUrlGroupBlacklists []AdUrlGroupBlacklist

func (a *AdUrlGroupBlacklists) FindAll() error {
	return BroadvidDB.Table("ad_url_group_blacklists").Find(&a).Error
}

type AdUrlGroupBlacklist struct {
	ID   int64    `form:"id"`
	Name string   `form:"name"`
	Urls []string `form:"url[]" sql:"-"`
}

func (a *AdUrlGroupBlacklist) Find(id string) error {
	err := BroadvidDB.Table("ad_url_group_blacklists").Find(&a, id).Error

	if err != nil {
		return err
	}

	var aus []AdUrlBlacklist

	err = BroadvidDB.Where("ad_url_group_id = ?", a.ID).Find(&aus).Error

	if err != nil {
		return err
	}

	for _, u := range aus {
		a.Urls = append(a.Urls, u.Url)
	}

	return nil
}

func (a *AdUrlGroupBlacklist) Save() error {
	err := BroadvidDB.Table("ad_url_group_blacklists").Save(&a).Error

	if err != nil {
		return err
	}

	var au AdUrlBlacklist

	err = BroadvidDB.Where("ad_url_group_id = ?", a.ID).Delete(&au).Error

	if err != nil {
		return err
	}

	for _, u := range a.Urls {
		var au AdUrlBlacklist
		au.Url = u
		au.AdUrlGroupID = a.ID
		BroadvidDB.Save(&au)
	}

	return nil
}

type AdUrlBlacklist struct {
	ID           int64  `form:"id"`
	Url          string `form:"url"`
	AdUrlGroupID int64  `form:"ad_url_group_id"`
}
