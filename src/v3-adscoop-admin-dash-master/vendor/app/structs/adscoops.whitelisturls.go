package structs

type AdscoopWhitelist struct {
	ID   uint     `form:"id"`
	Name string   `form:"name"`
	Urls []string `form:"url[]" sql:"-"`
}

type AdscoopWhitelistUrl struct {
	ID                 uint   `form:"id"`
	Url                string `form:"url"`
	AdscoopWhitelistId uint   `form:"whitelist_id"`
}

type AdscoopWhitelists []AdscoopWhitelist

func (a *AdscoopWhitelists) FindAll() error {
	return AdscoopsDB.Table("adscoop_whitelists").Find(&a).Error
}

func (a *AdscoopWhitelist) Find(id string) error {
	err := AdscoopsDB.Find(&a, id).Error

	if err != nil {
		return err
	}

	var asus []AdscoopWhitelistUrl

	err = AdscoopsDB.Where("adscoop_whitelist_id = ?", a.ID).Find(&asus).Error

	if err != nil {
		return err
	}

	for _, u := range asus {
		a.Urls = append(a.Urls, u.Url)
	}
	return nil
}

func (a *AdscoopWhitelist) Save() error {
	err := AdscoopsDB.Save(&a).Error

	if err != nil {
		return err
	}

	var asu AdscoopWhitelistUrl
	err = AdscoopsDB.Where("adscoop_whitelist_id = ?", a.ID).Delete(&asu).Error

	if err != nil {
		return err
	}

	for _, u := range a.Urls {
		var asu AdscoopWhitelistUrl
		asu.Url = u
		asu.AdscoopWhitelistId = a.ID
		if err := AdscoopsDB.Save(&asu).Error; err != nil {
			return err
		}
	}

	return nil
}
