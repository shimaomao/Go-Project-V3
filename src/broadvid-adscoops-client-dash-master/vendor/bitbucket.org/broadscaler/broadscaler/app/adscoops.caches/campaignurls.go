package adscoopsCaches

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/structs"
)

func LoadAdscoopCampaignUrls(id uint) (retAsu []structs.CampaignUrl) {
	var asu []structs.CampaignUrl
	location, _ := time.LoadLocation("America/Los_Angeles")

	today := time.Now()

	today = today.In(location)

	today = time.Date(today.Year(),
		today.Month(),
		today.Day(), 0, 0, 0, 0, location)

	today = today.In(time.UTC)
	lastTwoHours := time.Now().In(location).Add(-2 * time.Hour).In(time.UTC)

	var campaign structs.Campaign

	if err := campaign.Find(fmt.Sprintf("%v", id)); err != nil {
		log.Errorf("Cannot load campaign: %s", err)
		return nil
	}

	AdscoopsDB.Select("adscoop_urls.*").
		Joins(`LEFT OUTER JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.id
			   JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id`).
		Where(`adscoop_urls.campaign_id = ? AND (adscoop_urls.weight != 0 OR adscoop_campaigns.type != 2)
			AND (((adscoop_trackings.timeslice >= ?) = 0) OR (adscoop_trackings.timeslice IS NULL || adscoop_trackings.timeslice >= ?))`, id, lastTwoHours, lastTwoHours).
		Group("adscoop_urls.id").
		Order("SUM(adscoop_trackings.count) / adscoop_urls.weight ASC").
		Find(&asu)

	for _, y := range asu {
		for n := 0; uint(n) < y.Weight; n += 1 {
			if campaign.AppendRc {
				if strings.Contains(y.Url, "?") {
					y.Url = y.Url + "&rc=true"
				} else {
					y.Url = y.Url + "?rc=true"
				}
			}
			retAsu = append(retAsu, y)
		}
	}

	b, err := json.Marshal(retAsu)

	if err == nil {
		rp := RedisPool.Get()
		defer rp.Close()

		rp.Do("SET", fmt.Sprintf(URLS_KEY, id), b)
	}
	return
}
