package adscoopUtils

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
	"github.com/jinzhu/gorm"
)

type UtilManager struct {
	db *gorm.DB
}

func NewUtilManager(db *gorm.DB) *UtilManager {
	var um UtilManager
	um.db = db
	return &um
}

func (u *UtilManager) IngestXml(campaignId uint) error {
	var campaign structs.AdscoopCampaign
	var asu structs.AdscoopCampaignUrl
	u.db.Find(&campaign, campaignId)

	if campaign.ID == 0 {
		return errors.New("Campaign not found")
	}

	u.db.Where("campaign_id = ?", campaign.ID).Find(&asu)

	if asu.ID == 0 {
		return errors.New("No URL found for this campaign")
	}

	client := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}

	req, err := http.NewRequest("GET", asu.Url, nil)

	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var retXML structs.RssXml

	err = xml.Unmarshal(body, &retXML)

	if err != nil {
		return err
	}

	if len(retXML.Channel.Items) == 0 {
		return errors.New("Not enough items to ingest")
	}

	u.db.Where("campaign_id = ? AND weight != 0", campaign.ID).Delete(structs.AdscoopCampaignUrl{})

	for _, item := range retXML.Channel.Items {
		var url structs.AdscoopCampaignUrl

		url.Url = item.Link
		url.Title = item.Title
		url.Weight = 1
		url.AdscoopCampaignID = campaign.ID

		info := u.db.Save(&url)

		if info.Error != nil {
			u.db.Where("campaign_id = ? AND url = ?", campaign.ID, item.Link).Find(&url)

			if url.ID == 0 {
				_, err := u.db.DB().Exec("UPDATE adscoop_urls SET deleted_at = '0000-00-00', weight = 1, title = ? WHERE campaign_id = ? AND url = ?", item.Title, campaign.ID, item.Link)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
