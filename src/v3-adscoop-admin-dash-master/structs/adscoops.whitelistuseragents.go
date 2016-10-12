package structs

type AdscoopWhitelistUseragentGroup struct {
	ID         uint     `form:"id"`
	Name       string   `form:"name"`
	Useragents []string `form:"ua[]" sql:"-"`
}

func (w AdscoopWhitelistUseragentGroup) TableName() string {
	return "adscoop_whitelist_useragent_groups"
}

type AdscoopWhitelistUseragent struct {
	ID                               uint   `form:"id"`
	Useragent                        string `form:"ua"`
	AdscoopWhitelistUseragentGroupId uint   `form:"whitelist_useragent_group"`
}

func (w AdscoopWhitelistUseragent) TableName() string {
	return "adscoop_whitelist_useragents"
}

type AdscoopWhitelistUseragentGroups []AdscoopWhitelistUseragentGroup

func (a *AdscoopWhitelistUseragentGroups) FindAll() error {
	return AdscoopsDB.Table("adscoop_whitelist_useragent_groups").Find(&a).Error
}

func (a *AdscoopWhitelistUseragentGroup) Find(id string) error {
	err := AdscoopsDB.Find(&a, id).Error

	if err != nil {
		return err
	}

	var asu []AdscoopWhitelistUseragent

	err = AdscoopsDB.Where("adscoop_whitelist_useragent_group_id = ?", a.ID).
		Find(&asu).Error

	if err != nil {
		return err
	}

	for _, ua := range asu {
		a.Useragents = append(a.Useragents, ua.Useragent)
	}

	return nil
}

func (a *AdscoopWhitelistUseragentGroup) Save() error {
	err := AdscoopsDB.Save(&a).Error

	if err != nil {
		return err
	}

	var asu AdscoopWhitelistUseragent
	err = AdscoopsDB.Where("adscoop_whitelist_useragent_group_id = ?", a.ID).Delete(&asu).Error

	if err != nil {
		return err
	}

	for _, u := range a.Useragents {
		var asu AdscoopWhitelistUseragent
		asu.Useragent = u
		asu.AdscoopWhitelistUseragentGroupId = a.ID
		if err := AdscoopsDB.Save(&asu).Error; err != nil {
			return err
		}
	}

	return nil
}
