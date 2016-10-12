package structs

type AdCountryGroups []AdCountryGroup

func (a *AdCountryGroups) FindAll() error {
	return BroadvidDB.Table("ad_country_groups").Find(&a).Error
}

type AdCountryGroup struct {
	ID        int64    `form:"id"`
	Name      string   `form:"name"`
	Countries []string `form:"country[]" sql:"-"`
}

func (a *AdCountryGroup) Find(id string) error {
	err := BroadvidDB.Table("ad_country_groups").Find(&a, id).Error

	if err != nil {
		return err
	}

	var acs []AdCountry

	err = BroadvidDB.Where("ad_country_group_id = ?", a.ID).Find(&acs).Error

	if err != nil {
		return err
	}

	for _, u := range acs {
		a.Countries = append(a.Countries, u.CountryCode)
	}

	return nil
}

func (a *AdCountryGroup) Save() error {
	err := BroadvidDB.Table("ad_country_groups").Save(&a).Error

	if err != nil {
		return err
	}

	var ac AdCountry
	err = BroadvidDB.Where("ad_country_group_id = ?", a.ID).Delete(&ac).Error

	if err != nil {
		return err
	}

	for _, u := range a.Countries {
		var au AdCountry
		au.CountryCode = u
		au.AdCountryGroupID = a.ID
		BroadvidDB.Save(&au)
	}

	return nil
}

type AdCountry struct {
	ID               int64  `form:"id"`
	CountryCode      string `form:"country"`
	AdCountryGroupID int64  `form:"ad_country_group_id"`
}
