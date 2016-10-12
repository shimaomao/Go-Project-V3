package structs

type AdUaGroups []AdUaGroup

func (a *AdUaGroups) FindAll() error {
	return BroadvidDB.Table("ad_ua_groups").Find(&a).Error
}

type AdUaGroup struct {
	ID   int64    `form:"id"`
	Name string   `form:"name"`
	Uas  []string `form:"ua[]" sql:"-"`
}

func (a *AdUaGroup) Find(id string) error {
	err := BroadvidDB.Table("ad_ua_groups").Find(&a, id).Error

	if err != nil {
		return err
	}

	var aus []AdUa

	err = BroadvidDB.Where("ad_ua_group_id = ?", a.ID).Find(&aus).Error

	if err != nil {
		return err
	}

	for _, u := range aus {
		a.Uas = append(a.Uas, u.UserAgent)
	}
	return nil
}

func (a *AdUaGroup) Save() error {
	err := BroadvidDB.Table("ad_ua_groups").Save(&a).Error

	if err != nil {
		return err
	}

	var au AdUa
	err = BroadvidDB.Where("ad_ua_group_id = ?", a.ID).Delete(&au).Error

	if err != nil {
		return err
	}

	for _, u := range a.Uas {
		var au AdUa
		au.UserAgent = u
		au.AdUaGroupID = a.ID
		BroadvidDB.Save(&au)
	}

	return nil
}

type AdUa struct {
	ID          int64  `form:"id"`
	UserAgent   string `form:"user_agent"`
	AdUaGroupID int64  `form:"ad_ua_group_id"`
}
