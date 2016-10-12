package adscoopsCaches

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/broadscaler/broadscaler/app/structs"
	log "github.com/Sirupsen/logrus"
)

func updateCampaignUrlCaches() {
	var campaigns structs.Campaigns
	err := campaigns.GetRecent()

	if err != nil {
		log.Errorf("Error getting recent campaigns: %s", err)
		return
	}

	for _, c := range campaigns {
		log.Infof("Getting urls for campaign: %v", c.ID)
		LoadAdscoopCampaignUrls(c.ID)
	}
}

func UpdateCampaign(campaignID uint) {
	rc := RedisPool.Get()

	defer rc.Close()

	LoadAdscoopCampaignUrls(campaignID)

	var psdata structs.PubSubMessage
	psdata.Event = AdscoopPSEventCampaignUrls
	psdata.Data.IDInt = campaignID

	rc.Do("PUBLISH", AdscoopRedisPubSubChannel, psdata.JSONify())

	var campaignRedirs []structs.Redirect
	AdscoopsDB.Joins(`
					JOIN adscoop_redirect_campaigns ON adscoop_redirect_campaigns.redirect_id = adscoop_redirects.id
					`).Where(`adscoop_redirect_campaigns.campaign_id = ?`, campaignID).Find(&campaignRedirs)

	for _, r := range campaignRedirs {
		LoadAdscoopRedirs(r.Hash)
		LoadAdscoopRedirsCampaign(campaignID)

		var psdata structs.PubSubMessage
		psdata.Event = AdscoopPSEventRedirect
		psdata.Data.IDInt = r.ID
		psdata.Data.IDString = r.Hash

		rc.Do("PUBLISH", AdscoopRedisPubSubChannel, psdata.JSONify())
	}
}

func LoadAdscoopRedirsCampaign(id uint) (asc []structs.Campaign) {

	var redir structs.Redirect
	AdscoopsDB.Find(&redir, id)

	asc, errm := redir.GetActiveCampaigns()

	if errm != nil {
		log.Errorf("campaign mysql err: %s", errm)
	}

	b, err := json.Marshal(asc)

	if err == nil {
		rp := RedisPool.Get()
		defer rp.Close()

		_, err := rp.Do("SET", fmt.Sprintf(CAMPAIGNS_KEY, id), b)

		if err != nil {
			log.Errorf("campaign redis err: %s", err)
		}
	} else {
		log.Errorf("campaign json err: %s", err)
	}
	return
}
