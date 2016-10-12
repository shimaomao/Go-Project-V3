package structs

func NewRedirect() (*AdscoopRedirect, error) {
	return &AdscoopRedirect{}, nil
}

func (m *RedirectManager) LoadRedirectByHash(hash string) (*AdscoopRedirect, error) {
	var redir AdscoopRedirect
	err := m.DB.Where("hash = ?", hash).Find(&redir).Error

	return &redir, err
}

func NewCampaign() (*AdscoopCampaign, error) {
	return &AdscoopCampaign{}, nil
}

func NewCampaignUrl() (*AdscoopCampaignUrl, error) {
	return &AdscoopCampaignUrl{}, nil
}

func NewRedirectManager() *RedirectManager {
	return &RedirectManager{}
}

func NewCampaignManager() *CampaignManager {
	return &CampaignManager{}
}
