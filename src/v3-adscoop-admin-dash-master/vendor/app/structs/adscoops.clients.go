package structs

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"app/billing"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"

)

// Clients

type Client struct {
	gorm.Model
	Name                        string `form:"name"`
	Paused                      bool   `form:"paused"`
	DefaultCpc                  string `sql:"type:decimal(11,4)" form:"default_cpc"`
	DailyImpsLimit              string `form:"daily_imps_limit"`
	ApprovedToCharge            bool   `form:"approved_to_charge"`
	ChargeAmount                string `sql:"type:decimal(11,2)" form:"charge_amount"`
	ExpirationWarning           bool
	ExpirationNotice            bool
	HourlyReporting             bool `form:"hourly_reporting"`
	BreakDownReportByHour       bool `form:"break_down_report_by_hour"`
	EnableReporting             bool `form:"enable_reporting"`
	EnhancedReporting           bool `form:"enhanced_reporting"`
	EnableReportAccountBalance  bool `form:"enable_report_account_balance"`
	StripeToken                 string
	ReportLastSent              time.Time `form:"report_last_sent"`
	Emails                      []string  `form:"email[]" sql:"-"`
	CampaignEmails              []string  `form:"campaign_email[]" sql:"-"`
	Transactions                []ClientTransaction
	EnableClientLogin           bool `form:"enable_client_login"`
	ClientSchedulesPendApproval bool `form:"client_schedules_pend_approval"`
	ShowMtdSpendInReport        bool `form:"show_mtd_spend_in_report"`
	InGoodStanding              bool `form:"in_good_standing"`
	HideFromDash                bool
	LoadStats                   TrackingRows `sql:"-"`
	EngagementStats             TrackingRows `sql:"-"`
	ImpressionStats             TrackingRows `sql:"-"`
	RequiredImps                uint         `sql:"-"`
	ReceivedImps                uint         `sql:"-"`
	TodaySpend                  float64      `sql:"-"`
	TotalSpent                  float64      `sql:"-"`
	TOSAvg                      uint         `sql:"-"`
	ChartColor                  string
	UserSettings                ClientUserSettings `sql:"-"`
}

type ClientUserSettings struct {
	CampaignSort string
	ShowInfo     bool
	ClientOrder  uint
}

type ClientCompact struct {
	ClientID    uint `gorm:"primary_key"`
	Count       uint
	Engagement  uint
	Loads       uint
	Revenue     float64
	CollectedAt time.Time
}

func (cc ClientCompact) TableName() string {
	return "adscoop_clients_compacts"
}

type Clients []Client

type ClientRedirectStats map[string][]uint

func (c *Clients) TableName() string {
	return "adscoop_clients"
}

func (c Client) RedirRealtimeStats() ClientRedirectStats {
	var redirStats = make(ClientRedirectStats)

	var retData []struct {
		RedirectID string `gorm:"column:redirect_id"`
		Count      uint   `gorm:"column:count"`
	}

	currentTime := getNow()

	if err := AdscoopsDB.Raw(`
		SELECT CONCAT_WS(" | ", adscoop_redirects.name,adscoop_redirects.id) as redirect_id, SUM(temp_stats.count) as count FROM temp_stats
		JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
		JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
		JOIN adscoop_redirects ON adscoop_redirects.id = temp_stats.redirect_id
		WHERE timeslice >= ?
		AND adscoop_campaigns.client_id = ?
		GROUP BY redirect_id, timeslice
		ORDER BY redirect_id, timeslice
		`, currentTime.Add(-30*time.Minute).UTC(), c.ID).
		Scan(&retData).Error; err != nil {
		log.Errorf("Cannot grab stats: %s", err)

		return nil
	}

	for _, rd := range retData {
		redirStats[rd.RedirectID] = append(redirStats[rd.RedirectID], rd.Count)
	}

	return redirStats
}

func (c Client) AssociatedRedirects() (redirects []Redirect) {

	currentDay := getDayStart()

	if err := AdscoopsDB.Raw(`
		SELECT adscoop_redirects.* FROM temp_stats
		JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
		JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
		JOIN adscoop_redirects ON adscoop_redirects.id = temp_stats.redirect_id
		WHERE timeslice >= ?
		AND adscoop_campaigns.client_id = ?
		GROUP BY adscoop_redirects.id
		`, currentDay.UTC(), c.ID).
		Scan(&redirects).Error; err != nil {
		log.Errorf("Cannot grab redirects: %s", err)
		return
	}

	return
}

func (c Client) CompactOldData() error {
	today := now.New(time.Now().UTC()).BeginningOfDay()

	var retData ClientCompact

	if err := AdscoopsDB.Raw(`SELECT SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
					 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
					 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
				END) * adscoop_trackings.cpc) as revenue,
				SUM(count) as count,
				SUM(engagement) as engagement,
				SUM(adscoop_trackings.load) as loads
FROM adscoop_campaigns
JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID WHERE adscoop_campaigns.client_id = ? AND adscoop_trackings.timeslice < ?`, c.ID, today).Scan(&retData).Error; err != nil {
		return err
	}

	retData.ClientID = c.ID
	retData.CollectedAt = today

	tx := AdscoopsDB.Begin()

	if err := tx.Delete(ClientCompact{}, "client_id = ?", c.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&retData).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}

func (c *Client) Save() error {
	var asc Client
	if c.ID != 0 {
		AdscoopsDB.Find(&asc, c.ID)
		c.StripeToken = asc.StripeToken
	}
	err := AdscoopsDB.Save(&c).Error

	if err != nil {
		return err
	}

	AdscoopsDB.Where("client_id = ?", c.ID).Delete(ClientEmail{})

	for _, y := range c.Emails {
		var ase ClientEmail
		ase.Email = y
		ase.ClientID = c.ID

		info := AdscoopsDB.Save(&ase)

		if info.Error != nil {
			AdscoopsDB.Where("client_id = ? AND email = ?", c.ID, y).Find(&ase)

			if ase.ID == 0 {
				AdscoopsDB.DB().Exec("UPDATE adscoop_client_emails SET deleted_at = '0000-00-00' WHERE client_id = ? AND email = ?", c.ID, y)
			}
		}
	}

	AdscoopsDB.Where("client_id = ?", c.ID).Delete(ClientCampaignEmail{})

	for _, y := range c.CampaignEmails {
		var ase ClientCampaignEmail
		ase.Email = y
		ase.ClientID = c.ID

		info := AdscoopsDB.Save(&ase)

		if info.Error != nil {
			AdscoopsDB.Where("client_id = ? AND email = ?", c.ID, y).Find(&ase)

			if ase.ID == 0 {
				AdscoopsDB.DB().Exec("UPDATE adscoop_client_campaign_emails SET deleted_at = '0000-00-00' WHERE client_id = ? AND email = ?", c.ID, y)
			}
		}
	}

	if asc.Paused != c.Paused {
		log.Infof("Pause state changed, so we need to push updates across the redirs associated")

		var redirects []Redirect

		AdscoopsDB.LogMode(true)

		AdscoopsDB.Select("adscoop_redirects.*").Where(`adscoop_campaigns.client_id = ?
		AND (adscoop_campaigns.deleted_at IS NULL OR adscoop_campaigns.deleted_at <= '0001-01-02')
		AND (adscoop_redirect_campaigns.deleted_at IS NULL OR adscoop_redirect_campaigns.deleted_at <= '0001-01-02')`, c.ID).
			Joins(`JOIN adscoop_redirect_campaigns ON adscoop_redirect_campaigns.redirect_id = adscoop_redirects.id
		JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_redirect_campaigns.campaign_id`).
			Group("adscoop_redirects.id").
			Find(&redirects)

		AdscoopsDB.LogMode(false)

		for _, r := range redirects {
			r.Save()
		}
	}

	return nil
}

func (c *Client) CreateAPH() (aph AdscoopPaymentHash, err error) {
	if err = AdscoopsDB.Where("client_id = ?", c.ID).Delete(AdscoopPaymentHash{}).Error; err != nil {
		return
	}

	aph.ClientID = c.ID
	aph.Hash = randSeq(64)

	err = AdscoopsDB.Save(&aph).Error
	return
}

func (c *Client) Find(id string) error {
	err := AdscoopsDB.Where("id = ?", id).Find(&c).Error

	if err != nil {
		return err
	}

	var es []ClientEmail
	err = AdscoopsDB.Where("client_id = ?", c.ID).Find(&es).Error

	if err != nil {
		return err
	}

	for _, y := range es {
		c.Emails = append(c.Emails, y.Email)
	}

	var esc []ClientCampaignEmail
	err = AdscoopsDB.Where("client_id = ?", c.ID).Find(&esc).Error

	if err != nil {
		return err
	}

	for _, y := range esc {
		c.CampaignEmails = append(c.CampaignEmails, y.Email)
	}

	err = AdscoopsDB.Where("client_id = ?", c.ID).Order("id DESC").Find(&c.Transactions).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Charge(amount uint) (err error) {
	chargeID, err := billing.Charge(amount, c.StripeToken, c.Name)

	if err != nil {
		return
	}

	var asct ClientTransaction
	asct.ClientID = c.ID
	asct.AmountCharged = fmt.Sprintf("%v", amount)
	asct.TransactionId = chargeID
	asct.Successful = true
	asct.Attempts = 1

	return AdscoopsDB.Save(&asct).Error
}

func (c *Client) ManualCharge(amount uint) (err error) {
	var asct ClientTransaction
	asct.ClientID = c.ID
	asct.AmountCharged = fmt.Sprintf("%v", amount)
	asct.TransactionId = "MANUAL"
	asct.Successful = true
	asct.Attempts = 1

	return AdscoopsDB.Save(&asct).Error
}

func (c *Clients) FindAll() error {
	return AdscoopsDB.Table("adscoop_clients").Find(&c).Error
}

func (c *Clients) FindVisible(userID uint) error {

	var getClients []Client
	if err := AdscoopsDB.Table("adscoop_clients").Where("hide_from_dash = 0 ").Find(&getClients).Error; err != nil {
		return err
	}

	for _, gc := range getClients {
		if err := gc.getStats(); err != nil {
			log.Errorf("Not returning client on getStats: %v, err: %s", gc.ID, err)
			continue
		}


		if err := gc.getSpendData(); err != nil {
			log.Errorf("Not returning client on getSpendData: %v, err: %s", gc.ID, err)
			continue
		}

		if err := gc.getUserSettings(userID); err != nil {
			log.Errorf("Could not load user settings: %v, err :%s", gc.ID, err)
		}
		*c = append(*c, gc)
	}

	return nil
}

func (c *Client) getSpendData() (err error) {

	today := getDayStart()

	var retData RetRevData
	var revDataToday RetRevData

	currentHour := getHourStart()

	beginningOfMillenium := time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC)

	prevRevKey := fmt.Sprintf("client_rev_%v_%s", c.ID, currentHour)

	prevRev, found := Cache.Get(prevRevKey)

	if !found {

		query := revQuery(c.ID, beginningOfMillenium, currentHour.UTC())

		err = query.Find(&revDataToday).Error
		if err != nil {
			return
		}

		Cache.Set(prevRevKey, revDataToday, time.Hour)
	} else {

		revDataToday = prevRev.(RetRevData)
	}


	if err := revQuery(c.ID, currentHour.UTC(), today.UTC().Add(24*time.Hour)).Scan(&retData).Error; err != nil {
		return err
	}

	retData.TotalSpend += revDataToday.TotalSpend

	var totalCharged float64

	AdscoopsDB.Select(`SUM(amount_charged) as amount_charged`).Table(`adscoop_client_transactions`).Where("client_id = ?", c.ID).Row().Scan(&totalCharged)

	c.TotalSpent = totalCharged - retData.TotalSpend

	prevTOSKey := fmt.Sprintf("client_tos_%v_%s", c.ID, currentHour)

	prevTOS, found := Cache.Get(prevTOSKey)

	var retDataToday RetTOSData

	var retDataHour RetTOSData

	if !found {
		query := tosQuery(c.ID, today.UTC(), currentHour.UTC())

		err = query.Find(&retDataToday).Error
		if err != nil {
			return
		}

		Cache.Set(prevTOSKey, retDataToday, time.Hour)
	} else {
		retDataToday = prevTOS.(RetTOSData)
	}

	if err = tosQuery(c.ID, currentHour.UTC(), currentHour.UTC().Add(24*time.Hour)).Find(&retDataHour).Error; err != nil {
		return
	}

	c.TodaySpend = retDataToday.TotalSpend + retDataHour.TotalSpend
	c.TOSAvg = ((retDataToday.Tos * uint(currentHour.Hour())) + retDataHour.Tos) / uint(currentHour.Hour()+1)

	return
}

func (c *Client) getStats() error {
	var err error
	var impStatsCount, engagementStatsCount, loadStatsCount uint
	if impStatsCount, err = c.ImpressionStats.ForClient(c.ID, 0); err != nil {
		return err
	}

	if engagementStatsCount, err = c.EngagementStats.ForClient(c.ID, 1); err != nil {
		return err
	}

	if loadStatsCount, err = c.LoadStats.ForClient(c.ID, 2); err != nil {
		return err
	}

	c.ReceivedImps = impStatsCount + engagementStatsCount + loadStatsCount
	if err = c.GetRequiredImps(); err != nil {
		return err
	}


	return nil
}

func (c *Client) getUserSettings(userID uint) error {
	var userSettingsCampaignSort UserAdscoopsClientSetting
	AdscoopsDB.Where("user_id = ? AND client_id = ?", userID, c.ID).Find(&userSettingsCampaignSort)
	c.UserSettings.CampaignSort = userSettingsCampaignSort.CampaignSort
	c.UserSettings.ShowInfo = userSettingsCampaignSort.ShowInfo
	c.UserSettings.ClientOrder = userSettingsCampaignSort.ClientOrder
	return nil
}

func (c *Client) GetRequiredImps() error {
	var retData struct {
		RequiredImps uint
	}

	location, _ := time.LoadLocation("America/Los_Angeles")

	now := time.Now()
	now = now.In(location)

	err := AdscoopsDB.Select("SUM(adscoop_campaigns.daily_imps_limit) as required_imps").Table("adscoop_campaigns").
		Where(`
			(adscoop_campaigns.deleted_at IS NULL or adscoop_campaigns.deleted_at <= '0001-01-02')
			AND (adscoop_campaigns.paused = 0)
			AND (adscoop_campaigns.enable_start_stop_times = 0 OR ((adscoop_campaigns.disable_start_time = 1 OR adscoop_campaigns.start_datetime <= ?) AND (adscoop_campaigns.disable_end_time = 1 OR adscoop_campaigns.end_datetime > ?)))
			AND adscoop_campaigns.client_id = ?
			`, now.String(), now.String(), c.ID).Find(&retData).Error

	if err != nil {
		return err
	}

	c.RequiredImps = retData.RequiredImps

	return nil
}

func (c *Client) TableName() string {
	return "adscoop_clients"
}

type ClientTransaction struct {
	gorm.Model
	ClientID      uint   `form:"client_id"`
	AmountCharged string `sq:"type:decimal(11,2)" form:"amount_charged"`
	TransactionId string
	Successful    bool
	Attempts      uint
}

func (c *ClientTransaction) TableName() string {
	return "adscoop_client_transactions"
}

type ClientEmail struct {
	gorm.Model
	Email    string `sql:"unique_index:client_id"`
	ClientID uint   `sql:"unique_index:client_id"`
}

func (c *ClientEmail) TableName() string {
	return "adscoop_client_emails"
}

type ClientCampaignEmail struct {
	gorm.Model
	Email    string `sql:"unique_index:client_id"`
	ClientID uint   `sql:"unique_index:client_id"`
}

func (c *ClientCampaignEmail) TableName() string {
	return "adscoop_client_campaign_emails"
}

func tosQuery(clientID uint, timeslice time.Time, endTimeslice time.Time) *gorm.DB {
	return AdscoopsDB.Select(`
				   SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
						 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
						 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
					END) * adscoop_trackings.cpc) as total_spend,
					CASE WHEN SUM(time_on_site) / SUM(time_on_site_count)  > 0 THEN ROUND(SUM(time_on_site) / SUM(time_on_site_count)) else 0 END as tos`).
		Table("adscoop_campaigns").
		Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
				JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID`).
		Where("adscoop_campaigns.client_id = ? AND timeslice >= ? AND timeslice < ?", clientID, timeslice, endTimeslice)
}

func revQuery(clientID uint, timeslice time.Time, endTimeslice time.Time) *gorm.DB {

	return AdscoopsDB.LogMode(true).Raw(`SELECT SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
					 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
					 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
				END) * adscoop_trackings.cpc) as total_spend
FROM adscoop_campaigns
JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID WHERE adscoop_campaigns.client_id = ? AND adscoop_trackings.timeslice >= ? AND adscoop_trackings.timeslice < ?`, clientID, timeslice, endTimeslice)
}

type RetTOSData struct {
	TotalSpend float64
	Tos        uint
}

type RetRevData struct {
	TotalSpend float64
}
