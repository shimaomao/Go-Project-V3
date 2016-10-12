package structs

const ADSCOOPS_CAMPAIGN_GROUP = "adscoop_campaign_groups"

type CampaignGroup struct {
	ID   uint   `form:"id"`
	Name string `form:"name"`
}

func (c CampaignGroup) TableName() string {
	return ADSCOOPS_CAMPAIGN_GROUP
}
func (c *CampaignGroup) Save() error {
	return AdscoopsDB.Save(&c).Error
}

func (c *CampaignGroup) Find(id string) error {
	return AdscoopsDB.Table(ADSCOOPS_CAMPAIGN_GROUP).Find(&c, id).Error
}

type CampaignGroups []CampaignGroup

func (c *CampaignGroups) FindAll() error {
	return AdscoopsDB.Table(ADSCOOPS_CAMPAIGN_GROUP).Find(&c).Error
}
func (c CampaignGroups) TableName() string {
	return ADSCOOPS_CAMPAIGN_GROUP
}

type ClientCampaignGroupReads []ClientCampaignGroupRead

func (c *ClientCampaignGroupReads) Find(id string) error {
	return AdscoopsDB.
		Select("*").
		Table(ADSCOOPS_CAMPAIGN_GROUP).
		Joins("JOIN adscoop_client_campaign_groups ON adscoop_client_campaign_groups.campaign_group_id = adscoop_campaign_groups.id").
		Where("adscoop_client_campaign_groups.client_id = ?", id).
		Find(&c).Error
}

func (c ClientCampaignGroupReads) Save(id string) error {

	err := AdscoopsDB.Unscoped().Where("client_id = ?", id).Delete(&ClientCampaignGroup{}).Error
	if err != nil {
		return err
	}

	for _, c := range c {
		var rc ClientCampaignGroup
		rc.CampaignGroupID = c.CampaignGroupID
		rc.ClientID = c.ClientID

		err := AdscoopsDB.Save(&rc).Error

		if err != nil {
			return err
		}

	}

	return nil
}

type ClientCampaignGroupRead struct {
	ClientCampaignGroup
	Name string
}

type ClientCampaignGroup struct {
	ClientID        uint
	CampaignGroupID uint
}

func (c ClientCampaignGroup) TableName() string {
	return "adscoop_client_campaign_groups"
}

type RedirectCampaignGroupReads []RedirectCampaignGroupRead

func (c *RedirectCampaignGroupReads) Find(id string) error {
	return AdscoopsDB.
		Select("*").
		Table(ADSCOOPS_CAMPAIGN_GROUP).
		Joins("JOIN adscoop_redirect_campaign_groups ON adscoop_redirect_campaign_groups.campaign_group_id = adscoop_campaign_groups.id").
		Where("adscoop_redirect_campaign_groups.redirect_id = ?", id).
		Find(&c).Error
}

func (c RedirectCampaignGroupReads) Save(id string) error {
	err := AdscoopsDB.Unscoped().Where("redirect_id = ?", id).Delete(&RedirectCampaignGroup{}).Error

	if err != nil {
		return err
	}

	for _, c := range c {
		var rc RedirectCampaignGroup
		rc.CampaignGroupID = c.CampaignGroupID
		rc.RedirectID = c.RedirectID

		err := AdscoopsDB.Save(&rc).Error

		if err != nil {
			return err
		}
	}
	return nil
}

type RedirectCampaignGroupRead struct {
	RedirectCampaignGroup
	Name string
}

type RedirectCampaignGroup struct {
	RedirectID      uint
	CampaignGroupID uint
}

func (c RedirectCampaignGroup) TableName() string {
	return "adscoop_redirect_campaign_groups"
}
