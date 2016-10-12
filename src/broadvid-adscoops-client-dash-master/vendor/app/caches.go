package app

import (
	"fmt"
	"time"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
)

var campaignsCache = make(map[uint]CacheCampaign)

type CacheCampaign struct {
	LastUpdated       time.Time
	Campaigns         []CacheCampaignLayout
	ActiveCampaigns   []CacheCampaignLayout
	InactiveCampaigns []CacheCampaignLayout
	PausedCampaigns   []CacheCampaignLayout
}

type CacheCampaignLayout struct {
	structs.AdscoopCampaign
	Impressions    string
	Engagements    string
	Loads          string
	ImpressionsRaw uint
	EngagementsRaw uint
	LoadsRaw       uint
	Class          string
	DisplayImps    string
}

func setupCampaignCaches() {
	var clients []structs.AdscoopClient
	db.Where("enable_client_login = 1").Find(&clients)

	for _, c := range clients {
		createCampaignCache(c.ID)
	}
}

func createCampaignCache(clientid uint) {

	var wsData struct {
		Type      string
		Timestamp time.Time
	}
	tmpcc := CacheCampaign{}

	location, _ := time.LoadLocation("America/Los_Angeles")

	today := time.Now()

	today = today.In(location)

	wsData.Timestamp = today
	tmpcc.LastUpdated = today

	today = time.Date(today.Year(),
		today.Month(),
		today.Day(), 0, 0, 0, 0, location)

	today = today.In(time.UTC)

	db.LogMode(true).Select("adscoop_campaigns.id, adscoop_campaigns.name, adscoop_campaigns.campaign_group_weight, adscoop_campaigns.tracking_method, adscoop_campaigns.daily_imps_limit, adscoop_campaigns.paused, SUM(impressions.c) as impressions_raw, SUM(engagements.c) as engagements_raw, SUM(loads.c) as loads_raw, FORMAT(SUM(impressions.c),0) as impressions, FORMAT(SUM(engagements.c),0) as engagements, FORMAT(SUM(loads.c),0) as loads").
		Table("adscoop_campaigns").
		Joins(fmt.Sprintf(`LEFT OUTER JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
			LEFT OUTER JOIN (SELECT SUM(adscoop_trackings.count) as c, url_id FROM adscoop_trackings WHERE timeslice >= '%s' GROUP BY url_id) as impressions ON impressions.url_id = adscoop_urls.id
			LEFT OUTER JOIN (SELECT SUM(adscoop_trackings.engagement) as c, url_id FROM adscoop_trackings WHERE timeslice >= '%s' GROUP BY url_id) as engagements ON engagements.url_id = adscoop_urls.id
			LEFT OUTER JOIN (SELECT SUM(adscoop_trackings.load) as c, url_id FROM adscoop_trackings WHERE timeslice >= '%s' GROUP BY url_id) as loads ON loads.url_id = adscoop_urls.id`,
			today.Format(TimeLayout), today.Format(TimeLayout), today.Format(TimeLayout))).
		Where("adscoop_campaigns.client_id = ? AND inactive = 0", clientid).
		Order("paused asc").
		Group("adscoop_campaigns.id").
		Find(&tmpcc.Campaigns)

	for _, x := range tmpcc.Campaigns {
		x.Class = "success"

		if (x.TrackingMethod == 0 && x.DailyImpsLimit <= x.ImpressionsRaw) ||
			(x.TrackingMethod == 1 && x.DailyImpsLimit <= x.EngagementsRaw) ||
			(x.TrackingMethod == 2 && x.DailyImpsLimit <= x.LoadsRaw) {
			x.Class = "default"
			tmpcc.InactiveCampaigns = append(tmpcc.InactiveCampaigns, x)
			continue
		}

		if x.Paused {
			x.Class = "danger"
			tmpcc.PausedCampaigns = append(tmpcc.PausedCampaigns, x)
			continue
		}

		tmpcc.ActiveCampaigns = append(tmpcc.ActiveCampaigns, x)

	}

	campaignsCache[clientid] = tmpcc

	wsData.Type = "campaigns"

	broadcastJson(wsData, fmt.Sprintf("%v", clientid))
}
