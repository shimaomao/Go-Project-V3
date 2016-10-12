package structs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/patrickmn/go-cache"
	"app/shared-structs/structs"
	"github.com/jinzhu/gorm"
)

type Campaign struct {
	gorm.Model
	ClientID                    uint           `form:"client_id" sql:"index"`
	Cpc                         string         `sql:"type:decimal(11,6)" form:"cpc"`
	Name                        string         `form:"name"`
	Paused                      bool           `form:"paused"`
	DailyImpsLimit              uint           `form:"daily_imps_limit"`
	TrackingMethod              string         `form:"tracking_method"`
	Source                      string         `form:"source"`
	Type                        uint           `form:"type"`
	StartDatetimeEdit           string         `form:"start_datetime_edit" sql:"-"`
	EndDatetimeEdit             string         `form:"end_datetime_edit" sql:"-"`
	XmlType                     uint           `form:"xml_type"`
	XmlUrl                      string         `form:"xml_url" sql:"-"`
	PerformanceBasedPauseEnable bool           `form:"performance_based_pause_enable"`
	PerformanceBasedPauseResume bool           `form:"performance_based_pause_resume"`
	PerformanceBasedCompareA    string         `form:"performance_based_compare_a"`
	PerformanceBasedCompareB    string         `form:"performance_based_compare_b"`
	PerformanceBasedPercent     uint           `form:"performance_based_percent"`
	PerformanceBasedPauseQueued uint           `form:"performance_based_pause_queued"`
	PerformanceBasedNotifyOnly  bool           `form:"performance_based_notify_only"`
	EnableCampaignQualityCheck  bool           `form:"enable_campaign_quality_check"`
	StartDatetime               time.Time      `form:"-"`
	EndDatetime                 time.Time      `form:"-"`
	DisableStartTime            bool           `form:"disable_start_time"`
	DisableEndTime              bool           `form:"disable_end_time"`
	WeightVariance              uint           `form:"weight_variance"`
	EnableStartStopTimes        bool           `form:"enable_start_stop_times"`
	WeightsLastUpdated          time.Time      `form:"-"`
	Urls                        []*CampaignUrl `form:"-" sql:"-" json:"urls"`
	AllUrls                     []*CampaignUrl `form:"-" sql:"-" json:"all_urls"`
	ActiveWeights               []uint         `form:"weight[]" sql:"-"`
	Inactive                    bool           `form:"inactive"`
	EnableUnloadTracking        bool           `form:"enable_unload_tracking"`
	EnableSoftPause             bool
	IsRon                       bool
	Stats                       struct {
		DailyImps    float64
		DailyRevenue float64
	} `sql:"-"`
	CampaignGroupID     string
	CampaignGroupWeight string
	AppendRc            bool
}

func (c *Campaign) GetName() string {
	return c.Name
}

func (c Campaign) IsLimitReached() bool {
	today := getDayStart()
	todayUTC := today.UTC()
	var retData struct {
		IsTrue bool
	}
	var trackingType string

	switch c.TrackingMethod {
	case "0":
		trackingType = "count"
	case "1":
		trackingType = "engagement_count"
	case "2":
		trackingType = "load_count"
	}

	selectQuery := fmt.Sprintf(`SUM(temp_stats.%s) > %v as is_true`, trackingType, c.DailyImpsLimit)

	err := AdscoopsDB.
		Select(selectQuery).
		Table("temp_stats").
		Joins("JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id").
		Where("adscoop_urls.campaign_id = ? AND timeslice >= ?", c.ID, todayUTC).
		Find(&retData).Error

	if err != nil {
		log.Errorf("IsLimitReached error: %s", err)
		return true // going to force it to pause on error to yield to caution
	}
	return retData.IsTrue
}

type CampaignTodayCount struct {
	CampaignID uint
	Count      uint
}

func (c CampaignTodayCount) TableName() string {
	return "adscoop_campaign_today_count"
}

func (c Campaign) UpdateTodayCount() {
	today := getDayStart()
	todayUTC := today.UTC()
	var retData struct {
		Count uint
	}

	var countSelector = "SUM(count) as count"

	if c.TrackingMethod == "1" {
		countSelector = "SUM(engagement_count) as count"
	}
	if c.TrackingMethod == "2" {
		countSelector = "SUM(`temp_stats`.`load_count`) as count"
	}

	err := AdscoopsDB.
		Select(countSelector).
		Table("temp_stats").
		Joins(`JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id`).
		Where("temp_stats.timeslice >= ? AND adscoop_urls.campaign_id = ?", todayUTC.String(), c.ID).
		Find(&retData).
		Error

	if err != nil {
		log.Errorf("Error getting the day count for the campaign: %s", err)
		return
	}

	err = AdscoopsDB.Exec(`DELETE FROM adscoop_campaign_today_count WHERE campaign_id = ?`, c.ID).Error
	if err != nil {
		log.Errorf("Could not delete the campaign count for the day: %s", err)
		return
	}
	err = AdscoopsDB.Exec(`INSERT INTO adscoop_campaign_today_count (campaign_id, count) VALUES(?,?)`, c.ID, retData.Count).Error

	if err != nil {
		log.Errorf("Could not update the campaign count for the day: %s", err)
		return
	}
}

func (c Campaign) IsClientGoodStanding() bool {
	var cacheKey = fmt.Sprintf("isClientGoodStanding%v", c.ClientID)

	foo, found := Cache.Get(cacheKey)

	if found {
		return foo.(bool)
	}

	var client Client
	err := client.Find(fmt.Sprintf("%v", c.ClientID))

	if err != nil || client.ID == 0 {
		Cache.Set(cacheKey, false, cache.DefaultExpiration)
		log.Errorf("Error finding client: %s", err)
		return false
	}

	if client.Paused {
		Cache.Set(cacheKey, false, cache.DefaultExpiration)
		return false
	}

	if client.InGoodStanding {
		Cache.Set(cacheKey, true, cache.DefaultExpiration)
		return true
	}

	var clientCharge struct {
		Cost          float64
		AmountCharged float64
	}

	err = AdscoopsDB.
		Select(`SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
					 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
					 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
				END) * adscoop_trackings.cpc) as cost`).
		Table("adscoop_trackings").
		Joins(`JOIN adscoop_urls ON adscoop_urls.id = adscoop_trackings.url_id
		JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
		JOIN adscoop_clients ON adscoop_clients.id = adscoop_campaigns.client_id`).
		Where("adscoop_clients.id = ?", c.ClientID).Find(&clientCharge).Error

	if err != nil {
		log.Errorf("Getting client spend amount error: %s", err)
		return false
	}

	AdscoopsDB.Select(`SUM(amount_charged) as amount_charged`).
		Table(`adscoop_client_transactions`).
		Where(`client_id = ?`, c.ClientID).Find(&clientCharge)

	var retData = clientCharge.Cost < clientCharge.AmountCharged

	Cache.Set(cacheKey, retData, cache.DefaultExpiration)

	return retData
}

type CampaignUrl struct {
	gorm.Model
	AdscoopCampaignID uint `gorm:"column:campaign_id" sql:"unique_index:campaign_url"`
	Weight            uint
	Url               string `sql:"unique_index:campaign_url"`
	Title             string
}

func (a CampaignUrl) TableName() string {
	return "adscoop_urls"
}

type Campaigns []Campaign

func (c Campaigns) UpdateTodayCount() {
	today := getDayStart()
	todayUTC := today.UTC()

	if err := AdscoopsDB.Delete(&CampaignTodayCount{}).Error; err != nil {
		log.Errorf("Cannot clear out campaign today count table")
		return
	}

	var urls []CampaignUrl
	err := AdscoopsDB.Unscoped().
		Select("adscoop_urls.*").
		Table("adscoop_urls").
		Joins(`JOIN temp_stats ON temp_stats.url_id = adscoop_urls.id`).
		Where("temp_stats.timeslice >= ?", todayUTC.String()).
		Group("adscoop_urls.campaign_id").
		Find(&urls).Error

	if err != nil {
		log.Printf("Cannot update today count: %s", err)
		return
	}

	for _, u := range urls {
		var campaign Campaign
		if err = campaign.Find(fmt.Sprintf("%v", u.AdscoopCampaignID)); err != nil {
			log.Errorf("Cannot find campaign: %s", err)
			continue
		}
		campaign.UpdateTodayCount()
	}
}

func (c *Campaigns) TableName() string {
	return "adscoop_campaigns"
}

func (c *Campaigns) GetRecent() error {
	// now := time.Now().UTC()
	// yesterday := now.Add(time.Duration(-24 * time.Hour))
	return AdscoopsDB.LogMode(true).Select("*").Table("adscoop_campaigns").Where("paused = 0").Find(&c).Error
	// return AdscoopsDB.Select("adscoop_campaigns.*").
	// 	Table("adscoop_campaigns").
	// 	Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
	// 			 JOIN temp_stats ON temp_stats.url_id = adscoop_urls.id`).
	// 	Where("temp_stats.timeslice >= ?", yesterday).Group("adscoop_urls.campaign_id").Find(&c).Error
}

func (c *Campaigns) FindAll() error {
	return AdscoopsDB.Table("adscoop_campaigns").Find(&c).Error
}

func (c *Campaigns) FindFromClient(clientID string) error {
	return AdscoopsDB.Table("adscoop_campaigns").Where("client_id = ? AND  (`adscoop_campaigns`.deleted_at IS NULL OR `adscoop_campaigns`.deleted_at <= '0001-01-02')", clientID).Find(&c).Error
}

func (c *Campaigns) FindExtraDetailsFromClient(clientID string) error {
	today := getDayStart()
	todayUTC := today.UTC()
	var campaigns []Campaign
	err := AdscoopsDB.Select("adscoop_campaigns.*").Table("adscoop_campaigns").Joins(`
		JOIN adscoop_urls ON adscoop_campaigns.id = adscoop_urls.campaign_id
		JOIN temp_stats ON adscoop_urls.id = temp_stats.url_id
		`).Where(`
			adscoop_campaigns.client_id = ?
			AND timeslice >= ?`, clientID, todayUTC.String()).Group(`adscoop_campaigns.id`).Order(`SUM(temp_stats.count) DESC`).Find(&campaigns).Error

	if err != nil {
		return err
	}

	for _, campaign := range campaigns {
		campaign.GetDailyStats()
		*c = append(*c, campaign)
	}

	return nil
}

func (c Campaign) TableName() string {
	return "adscoop_campaigns"
}

func (c *Campaign) GetDailyStats() {
	var countType string
	switch c.TrackingMethod {
	case "0":
		countType = "count"
	case "1":
		countType = "engagement_count"
	case "2":
		countType = "load_count"
	}

	var retData CampaignDailyStats
	var retDataHour CampaignDailyStats

	today := getDayStart()
	currentHour := getHourStart()

	prevDailyStatsKey := fmt.Sprintf("campaign_daily_stats_%s_%v_%s", countType, c.ID, currentHour)

	prevDailyStats, found := Cache.Get(prevDailyStatsKey)

	if !found {
		campaignDailyStatsQuery(countType, c.ID, today.UTC(), currentHour.UTC()).Find(&retData)
		Cache.Set(prevDailyStatsKey, retData, time.Hour)
	} else {
		retData = prevDailyStats.(CampaignDailyStats)
	}

	campaignDailyStatsQuery(countType, c.ID, currentHour.UTC(), currentHour.Add(time.Hour).UTC()).Find(&retDataHour)

	c.Stats.DailyImps = retData.DailyImps + retDataHour.DailyImps
	c.Stats.DailyRevenue = retData.DailyRevenue + retDataHour.DailyRevenue
}

func (c *Campaign) GetTempDailyStats() error {
	var retData struct {
		Count float64
	}
	err := AdscoopsDB.Select("count").Table("adscoop_campaign_today_count").Where("campaign_id = ?", c.ID).Find(&retData).Error

	if err != nil {
		return err
	}

	c.Stats.DailyImps = retData.Count
	return nil
}

func (c *Campaign) Find(id string) error {
	err := AdscoopsDB.Find(&c, id).Error

	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("America/Los_Angeles")
	c.StartDatetime = time.Date(c.StartDatetime.Year(), c.StartDatetime.Month(),
		c.StartDatetime.Day(), c.StartDatetime.Hour(), c.StartDatetime.Minute(), 0, 0, loc)

	log.Println("endDateTime", c.EndDatetime.String())
	c.EndDatetime = time.Date(c.EndDatetime.Year(), c.EndDatetime.Month(),
		c.EndDatetime.Day(), c.EndDatetime.Hour(), c.EndDatetime.Minute(), 0, 0, loc)

	if c.EndDatetime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		c.EndDatetime = time.Now()
	}

	if c.StartDatetime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		c.StartDatetime = time.Now()
	}

	if c.Type == 0 {
		var us []*CampaignUrl
		err := AdscoopsDB.Where("campaign_id = ?", c.ID).Order("id desc").Find(&us).Error

		if err != nil {
			return err
		}

		for _, y := range us {
			c.Urls = append(c.Urls, y)
		}

		us = []*CampaignUrl{}
		err = AdscoopsDB.Unscoped().Where("deleted_at != '0000-00-00' AND campaign_id = ?", c.ID).Order("id desc").Find(&us).Error

		if err != nil {
			return err
		}

		for _, y := range us {
			c.AllUrls = append(c.AllUrls, y)
		}
	}

	if c.Type == 1 || c.Type == 2 {
		var asu CampaignUrl
		err := AdscoopsDB.Where("campaign_id = ? AND weight = 0", c.ID).Find(&asu).Error

		if err != nil {
			return err
		}
		c.XmlUrl = asu.Url
	}

	return nil
}

func (c *Campaign) IngestXml() error {
	var asu CampaignUrl

	if c.ID == 0 {
		return errors.New("Campaign not found")
	}

	AdscoopsDB.Where("campaign_id = ?", c.ID).Find(&asu)

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

	AdscoopsDB.Where("campaign_id = ? AND weight != 0", c.ID).Delete(CampaignUrl{})

	for _, item := range retXML.Channel.Items {
		var url CampaignUrl

		url.Url = item.Link
		url.Title = item.Title
		url.Weight = 1
		url.AdscoopCampaignID = c.ID

		info := AdscoopsDB.Save(&url)

		if info.Error != nil {
			AdscoopsDB.Where("campaign_id = ? AND url = ?", c.ID, item.Link).Find(&url)

			if url.ID == 0 {
				_, err := AdscoopsDB.DB().Exec("UPDATE adscoop_urls SET deleted_at = '0000-00-00', weight = 1, title = ? WHERE campaign_id = ? AND url = ?", item.Title, c.ID, item.Link)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Campaign) BasicSave() error {
	var extc Campaign
	err := extc.Find(fmt.Sprintf("%v", c.ID))

	if err != nil {
		return err
	}

	extc.Paused = c.Paused
	extc.DailyImpsLimit = c.DailyImpsLimit
	extc.Cpc = c.Cpc

	err = extc.Save()

	return err
}

func (c *Campaign) Save() error {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	c.StartDatetime = c.StartDatetime.In(loc)
	c.EndDatetime = c.EndDatetime.In(loc)

	c.StartDatetime = time.Date(c.StartDatetime.Year(), c.StartDatetime.Month(),
		c.StartDatetime.Day(), c.StartDatetime.Hour(), c.StartDatetime.Minute(), 0, 0, time.UTC)

	log.Println("endDateTime", c.EndDatetime.String())
	c.EndDatetime = time.Date(c.EndDatetime.Year(), c.EndDatetime.Month(),
		c.EndDatetime.Day(), c.EndDatetime.Hour(), c.EndDatetime.Minute(), 0, 0, time.UTC)

	c.PerformanceBasedPauseQueued = 1
	err := AdscoopsDB.Save(&c).Error

	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())

	if c.Type == 0 {
		AdscoopsDB.Where("campaign_id = ?", c.ID).Delete(CampaignUrl{})
		for _, y := range c.Urls {
			var asu CampaignUrl
			asu.Url = y.Url
			var urlWeight uint
			if c.WeightVariance != 0 {
				urlWeight = uint(rand.Intn(100-(100-int(c.WeightVariance))) + (100 - int(c.WeightVariance)))
			} else {
				urlWeight = y.Weight
			}
			asu.Weight = urlWeight
			asu.AdscoopCampaignID = c.ID
			info := AdscoopsDB.Save(&asu)

			if info.Error != nil {
				AdscoopsDB.Where("campaign_id = ? AND url = ?", c.ID, y.Url).Find(&asu)
				if asu.ID == 0 {
					_, err := AdscoopsDB.DB().Exec("UPDATE adscoop_urls SET deleted_at = '0000-00-00', weight = ? WHERE campaign_id = ? AND url = ?", urlWeight, c.ID, y.Url)
					log.Println("url update err", err)
				}
			}
		}
	}

	if c.Type == 1 || c.Type == 2 {
		var asu CampaignUrl
		log.Println("c", c.ID)
		AdscoopsDB.Where("campaign_id = ? AND weight = 0", c.ID).Find(&asu)
		log.Println("url", c.XmlUrl)
		log.Println("asu", asu.ID)
		asu.Url = c.XmlUrl
		asu.AdscoopCampaignID = c.ID
		asu.Weight = 0
		AdscoopsDB.Save(&asu)
		if c.Type == 2 {
			c.IngestXml()
		}
	}

	return nil

}

type CampaignSchedule struct {
	Campaign
	CampaignScheduleAddons
	MacroReplace string `sql:"-"`
	MacroFind    string `sql:"-"`
}

func (c *CampaignSchedule) Find(id string) error {
	return AdscoopsDB.Find(&c, id).Error
}

func (c *CampaignSchedule) Save() error {

	if c.Type == 0 {
		if len(c.Urls) == 0 {
			c.Paused = true
		}
	}

	loc, _ := time.LoadLocation("America/Los_Angeles")

	c.StartDatetime = c.StartDatetime.In(loc)
	c.EndDatetime = c.EndDatetime.In(loc)

	today := time.Date(c.ScheduleExecute.Year(),
		c.ScheduleExecute.Month(),
		c.ScheduleExecute.Day(),
		c.ScheduleExecute.Hour(),
		c.ScheduleExecute.Minute(), 0, 0, loc)

	today = today.In(time.UTC)

	c.ScheduleExecute = today

	c.PerformanceBasedPauseQueued = 1

	err := AdscoopsDB.Save(&c).Error
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())

	if c.MacroReplace != "" {
		if c.MacroFind == "" {
			c.MacroFind = "[REPLACE_ME]"
		}

		for x, y := range c.Urls {
			c.Urls[x].Url = strings.Replace(y.Url, c.MacroFind, c.MacroReplace, -1)
		}
	}

	if c.Type == 0 {
		AdscoopsDB.Where("campaign_id = ?", c.ID).Delete(CampaignScheduleUrl{})
		for x, y := range c.Urls {
			var asu CampaignScheduleUrl
			asu.Url = y.Url
			var urlWeight uint
			if c.WeightVariance != 0 {
				urlWeight = uint(rand.Intn(100-(100-int(c.WeightVariance))) + (100 - int(c.WeightVariance)))
			} else {
				urlWeight = c.ActiveWeights[x]
			}
			asu.Weight = urlWeight
			asu.AdscoopCampaignID = c.ID
			info := AdscoopsDB.Save(&asu)

			if info.Error != nil {
				AdscoopsDB.Where("campaign_id = ? AND url = ?", c.ID, y).Find(&asu)
				if asu.ID == 0 {
					_, err := AdscoopsDB.DB().Exec("UPDATE adscoop_scheduled_urls SET deleted_at = '0000-00-00', weight = ? WHERE campaign_id = ? AND url = ?", urlWeight, c.ID, y)
					log.Println("url update err", err)
				}
			}
		}
	}

	return nil

}

func (c *CampaignSchedule) TableName() string {
	return "adscoop_campaign_schedules"
}

type CampaignScheduleUrl struct {
	CampaignUrl
}

func (a CampaignScheduleUrl) TableName() string {
	return "adscoop_scheduled_urls"
}

type CampaignSchedules []CampaignScheduleRead

type CampaignScheduleRead struct {
	gorm.Model
	CampaignScheduleAddons
}

func (c *CampaignSchedules) FindAll(campaignID string) error {
	err := AdscoopsDB.Table("adscoop_campaign_schedules").
		Where("campaign_id = ?", campaignID).
		Find(&c).Error

	if err != nil {
		return err
	}

	for i := 0; i < len((*c)); i++ {
		cs := (*c)[i]

		pst := time.Date(cs.ScheduleExecute.Year(),
			cs.ScheduleExecute.Month(),
			cs.ScheduleExecute.Day(),
			cs.ScheduleExecute.Hour(),
			cs.ScheduleExecute.Minute(),
			0, 0, time.UTC)

		location, _ := time.LoadLocation("America/Los_Angeles")

		cs.ScheduleExecute = pst.In(location)

		(*c)[i] = cs
	}

	return nil
}

type CampaignScheduleAddons struct {
	ScheduleExecuteEdit string    `sql:"-" form:"schedule_execute_edit"`
	ScheduleExecute     time.Time `form:"-"`
	CampaignID          uint      `form:"campaign_id" sql:"index"`
	ScheduleLabel       string    `form:"schedule_label"`
	ScheduleQueued      bool      `form:"schedule_queued"`
	SchedulePending     bool      `form:"schedule_pending"`
}

type CampaignDailyStats struct {
	DailyImps    float64
	DailyRevenue float64
}

func campaignDailyStatsQuery(countType string, campaignID uint, timeslice time.Time, endTimeslice time.Time) *gorm.DB {
	return AdscoopsDB.Raw(fmt.Sprintf(`SELECT SUM(temp_stats.%s) as daily_imps, SUM(temp_stats.revenue) as daily_revenue FROM temp_stats
JOIN adscoop_urls on temp_stats.url_id = adscoop_urls.id
WHERE adscoop_urls.campaign_id = ?
AND timeslice >= ? AND timeslice < ?`, countType), campaignID, timeslice.String(), endTimeslice.String())
}
