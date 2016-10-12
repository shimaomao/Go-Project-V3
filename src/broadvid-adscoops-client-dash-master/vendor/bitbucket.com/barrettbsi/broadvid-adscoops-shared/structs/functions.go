package structs

import (
	"errors"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (m *RedirectManager) Save(redir *AdscoopRedirect) error {
	if redir.Name == "" {
		return errors.New("Redir name not set")
	}
	if redir.Hash == "" {
		redir.Hash = randSeq(16)
	}
	return m.DB.Save(&redir).Error
}

func (m *RedirectManager) AssociateCampaign(redirect *AdscoopRedirect, campaign *AdscoopCampaign, weight uint) error {
	var rc AdscoopRedirectCampaign
	rc.CampaignID = campaign.ID
	rc.RedirectID = redirect.ID

	return m.DB.Save(&rc).Error
}

func (m *RedirectManager) LoadRedirByHash(hash string) (*AdscoopRedirect, error) {
	var redir AdscoopRedirect
	err := m.DB.Where("hash = ?", hash).Find(&redir).Error
	return &redir, err
}

func (m *CampaignManager) Save(campaign *AdscoopCampaign) error {
	if campaign.Name == "" {
		return errors.New("Campaign name not set")
	}
	for _, u := range campaign.Urls {
		if u.Url == "" {
			return errors.New("Campaign URL is blank")
		}
		if u.Weight == 0 {
			return errors.New("URL weight cannot equal 0")
		}
	}
	return m.DB.Save(&campaign).Error
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
