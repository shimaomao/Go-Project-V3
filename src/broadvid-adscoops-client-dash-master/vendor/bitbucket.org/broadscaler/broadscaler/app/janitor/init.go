package janitor

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/emailer"
	"bitbucket.org/broadscaler/broadscaler/app/structs"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool

func Run() {
	structs.AdscoopsDB = AdscoopsDB
	for t := range time.Tick(time.Minute * 1) {
		if t.Minute() == 0 || t.Minute() == 30 {
			runDropoff()
		}
	}
}

func TestRun() {
	structs.AdscoopsDB = AdscoopsDB
	runDropoff()
}

func runDropoff() {
	var retData struct {
		Urls []struct {
			ClientName   string
			ClientID     string
			CampaignName string
			CampaignID   string
			Url          string
			UrlID        uint
			Count        string
			Loads        string
			Dropoff      string
		}
	}

	location, _ := time.LoadLocation("America/Los_Angeles")

	today := time.Now()

	today = today.In(location)

	today = time.Date(today.Year(),
		today.Month(),
		today.Day(), 0, 0, 0, 0, location)

	today = today.In(time.UTC)

	AdscoopsDB.Select(`adscoop_clients.name as client_name, adscoop_clients.id as client_id,
            adscoop_campaigns.name as campaign_name, adscoop_campaigns.id as campaign_id,
            adscoop_urls.url as url, adscoop_urls.id as url_id,
            FORMAT(SUM(count),0) as count,
            FORMAT(SUM(adscoop_trackings.load),0) as loads,
            FLOOR((SUM(adscoop_trackings.load) / SUM(count) * 100)) as dropoff`).
		Table(`adscoop_trackings`).
		Joins(`JOIN adscoop_urls ON adscoop_urls.id = adscoop_trackings.url_id
             JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
             JOIN adscoop_redirects ON adscoop_redirects.id = adscoop_trackings.redirect_id
             JOIN adscoop_clients ON adscoop_clients.id = adscoop_campaigns.client_id`).
		Where(`adscoop_trackings.timeslice >= ?
             AND (adscoop_urls.deleted_at IS NULL OR adscoop_urls.deleted_at <= '0001-01-02')
             AND (adscoop_campaigns.deleted_at IS NULL OR adscoop_campaigns.deleted_at <= '0001-01-02')
             AND (adscoop_campaigns.tracking_method != 0 OR adscoop_redirects.iframe = 2)`, today).
		Group(`url_id`).
		Having(`SUM(count) * .25 > SUM(adscoop_trackings.load) AND SUM(count) > 100`).
		Order(`SUM(adscoop_trackings.load) asc`).
		Find(&retData.Urls)

	for _, u := range retData.Urls {
		var client structs.Client
		var campaign structs.Campaign

		if err := client.Find(u.ClientID); err != nil {
			log.Printf("Error finding client: %s", err)
			continue
		}

		if err := campaign.Find(u.CampaignID); err != nil {
			log.Printf("Error finding campaign: %s", err)
			continue
		}

		title := "Url going to be paused based off performance"
		message := fmt.Sprintf(`
    Url going to be paused based off performance: %s
    Count: %s
    Load: %s
`, u.Url, u.Count, u.Loads)

		client.CampaignEmails = append(client.CampaignEmails, "daniel.aharonoff@broadscaler.com")
		client.CampaignEmails = append(client.CampaignEmails, "adops@adscoops.com")

		var e emailer.Emailer

		err := e.Send(title, message, client.CampaignEmails)

		if err != nil {
			log.Errorf("Could not send email: %s", err)
			continue
		}

		var uToPause structs.CampaignUrl
		err = AdscoopsDB.Delete(uToPause, u.UrlID).Error
		if err != nil {
			log.Errorf("Could not pause URL: %s", err)
		}

	}
}
