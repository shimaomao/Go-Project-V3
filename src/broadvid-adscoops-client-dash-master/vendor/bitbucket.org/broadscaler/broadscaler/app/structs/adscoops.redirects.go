package structs

import (
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/broadscaler/broadscaler/app/adonnetwork"

	log "github.com/Sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Redirects

type Redirect struct {
	gorm.Model
	Name                    string
	Hash                    string `sql:"unique_index"`
	Min                     uint
	Max                     uint
	Iframe                  uint
	RedirType               uint
	BapiScoring             uint
	LockWhitelistId         string
	LockUseragentId         string
	LockWhitelistReverse    bool
	LockUseragentReverse    bool
	Paused                  bool
	LockWhitelistUrls       []AdscoopWhitelistUrl       `sql:"-"`
	LockUseragents          []AdscoopWhitelistUseragent `sql:"-"`
	AutoRefresh             bool
	StripReferrer           bool
	ForceRefresh            bool
	StripQueryString        bool
	ForceHost               string
	AllowedQueryStrings     []string `sql:"-"`
	BbsiPath                string
	SortMethod              uint
	ForceHostString         string                 `sql:"-"`
	Campaigns               []RedirectCampaignRead `sql:"-"`
	ScoringRedirectEnabled  bool
	ScoringRedirectOverride string
	ScoringTimeout          uint
	HideFromDash            bool
	CampaignSource          string
	BustIframe              bool

	/* This is to send a notification e-mail that the threshold for pause/unpause has been triggered for a redir */
	NotifyCampaignThreshold  bool
	CampaignThresholdPercent float64
	IsPerformancePaused      bool

	/* This is to enable notificatione e-mail of budget changes */
	AdvertisingDailySpend uint
	AdvertisingBid        string

	/* Adon API Integration */
	AdvertisingCampaignID         string
	EnableAdvertisingPause        bool
	EnableAdvertisingSpendChange  bool
	EnableAdvertisingBidChange    bool
	EnableAdvertisingASAPSpending bool `gorm:"column:enable_advertising_asap_spending"`
}

func (r *Redirect) TableName() string {
	return "adscoop_redirects"
}

func (r *Redirect) Find(id string) error {
	return AdscoopsDB.Find(&r, id).Error
}

func (r *Redirect) Save() error {
	if r.ID == 0 {
		r.Hash = randSeq(16)
	} else {
		var rold Redirect
		err := AdscoopsDB.Find(&rold, r.ID).Error

		if err != nil || rold.ID == 0 {
			return err
		}

		r.Hash = rold.Hash
		r.CreatedAt = rold.CreatedAt
	}

	err := AdscoopsDB.Save(&r).Error

	if err != nil {
		return err
	}

	err = AdscoopsDB.
		Where("redirect_id = ?", r.ID).
		Delete(RedirectQuerystring{}).Error

	if err != nil {
		return err
	}

	for _, y := range r.AllowedQueryStrings {
		var asqs RedirectQuerystring
		asqs.QueryStringKey = y
		asqs.RedirectID = r.ID
		if err := asqs.Save(); err != nil {
			err := AdscoopsDB.
				Where("redirect_id = ? AND query_string_key = ?", r.ID, y).
				Find(&asqs).Error

			if err != nil || asqs.ID != 0 {
				return err
			}

			err = AdscoopsDB.Exec(`UPDATE adscoop_redirect_querystrings SET deleted_at = '0000-00-00' WHERE redirect_id = ? AND query_string_key = ?`, r.ID, y).Error

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Redirect) CheckDailyProgress() error {
	if r.NotifyCampaignThreshold { // Check to see if notifications should be sent.

		if err := AdscoopsDB.Find(&r, r.ID).Error; err != nil {
			log.Errorf("Cannot find redir: %s", err)
			return err
		}

		var redirCampaignLimit float64
		var redirCampaignPerformance float64

		campaigns, err := r.GetLiveCampaigns(false)

		if err != nil {
			log.Errorf("Cannot get live campaigns for redir ID: %v", r.ID)
			return err
		}

		for _, c := range campaigns {
			log.Infof("Campaign name: %s", c.Name)
			redirCampaignLimit += float64(c.DailyImpsLimit)

			if err := c.GetTempDailyStats(); err != nil {
				log.Errorf("Cannot get temp daily stats: %s", err)
				continue
			}

			log.Infof("Campaign stats: %v", c.Stats.DailyImps)

			redirCampaignPerformance += c.Stats.DailyImps
		}

		var redirPerformance float64 = redirCampaignPerformance / redirCampaignLimit

		log.Infof("Redir ID: %v", r.ID)
		log.Infof("Redir Name: %s", r.Name)
		log.Infof("Redir Campaign Performance: %v", redirCampaignPerformance)
		log.Infof("Redir Campaign Limit: %v", redirCampaignLimit)
		log.Infof("Redir Performance: %v", redirPerformance)
		log.Infof("CampaignThresholdPercent: %v", r.CampaignThresholdPercent)
		log.Infof("IsPerformancePaused: %v", r.IsPerformancePaused)
		log.Infof("redirPerformance*100 >= r.CampaignThresholdPercent: %v", redirPerformance*100 >= r.CampaignThresholdPercent)
		log.Infof("(redirPerformance*100 >= r.CampaignThresholdPercent && !r.IsPerformancePaused): %v", (redirPerformance*100 >= r.CampaignThresholdPercent && !r.IsPerformancePaused))
		log.Infof("(len(campaigns) == 0 && !r.IsPerformancePaused): %v", (len(campaigns) == 0 && !r.IsPerformancePaused))

		if (redirPerformance*100 >= r.CampaignThresholdPercent && !r.IsPerformancePaused) || (len(campaigns) == 0 && !r.IsPerformancePaused) {
			log.Infof("Going to pause redir: %v", r.ID)
			r.IsPerformancePaused = true
			if err := r.Save(); err != nil {
				return err
			}

			r.PauseLinksInCampaigns()

			title := fmt.Sprintf("%s is going to be paused because it has completed for the day", r.Name)
			message := fmt.Sprintf("%s is at %v percent completion for the day.", r.Name, redirPerformance*100)

			var users Users
			if err := users.MessageAdscoopsUsers(title, message); err != nil {
				return err
			}

			if r.EnableAdvertisingPause {
				var ac adonnetwork.Campaign

				if err := ac.Find(r.AdvertisingCampaignID); err != nil {
					log.Errorf("Cannot load advertising campaign: %s", err)
					return err
				}

				for i := 0; i < 4; i++ {
					if i == 3 {
						title := fmt.Sprintf("ERROR: %s is not going to be paused because of an API error", r.Name)
						message := fmt.Sprintf("ERROR: %s is at %v percent completion for the day, but the request to pause failed %v times", r.Name, redirPerformance*100, i)

						users.MessageAdscoopsUsers(title, message)
						return nil
					}
					if err := ac.Pause(); err != nil {
						log.Errorf("Cannot pause advertising campaign: %s", err)
					} else {
						break
					}

				}

			}
		}

		if redirPerformance*100 < r.CampaignThresholdPercent && r.IsPerformancePaused && len(campaigns) != 0 {
			log.Infof("Going to unpause redir: %v", r.ID)
			r.IsPerformancePaused = false
			if err := r.Save(); err != nil {
				return err
			}

			r.UnpauseLinksInCampaigns()

			title := fmt.Sprintf("%s is going to be unpaused because it has not completed for the day", r.Name)
			message := fmt.Sprintf("%s is at %v percent completion for the day.", r.Name, redirPerformance*100)

			var users Users
			if err := users.MessageAdscoopsUsers(title, message); err != nil {
				return err
			}

			if r.EnableAdvertisingPause {
				var ac adonnetwork.Campaign
				if err := ac.Find(r.AdvertisingCampaignID); err != nil {
					log.Errorf("Cannot load advertising campaign: %s", err)
					return err
				}

				for i := 0; i < 4; i++ {
					if i == 3 {
						title := fmt.Sprintf("ERROR: %s is not going to be unpaused because of an API error", r.Name)
						message := fmt.Sprintf("ERROR: %s is at %v percent completion for the day, but the request to unpause failed %v times", r.Name, redirPerformance*100, i)

						users.MessageAdscoopsUsers(title, message)
						return nil
					}
					if err := ac.Resume(); err != nil {
						log.Errorf("Cannot resume advertising campaign: %s", err)
					} else {
						break
					}
				}

			}
		}

		if r.EnableAdvertisingBidChange {
			log.Infof("Bid change is enabled for redir: %v", r.ID)

			var rp RedirectPerformance
			err := AdscoopsDB.Where("redirect_id = ?", r.ID).Find(&rp).Error

			if err != nil || rp.Bid != r.AdvertisingBid {
				rp.RedirectID = r.ID
				rp.Bid = r.AdvertisingBid
				err := AdscoopsDB.Save(&rp).Error
				if err != nil {
					log.Errorf("Could not save redirect performance: %s", err)
					return err
				}

				var ac adonnetwork.Campaign
				if err = ac.Find(r.AdvertisingCampaignID); err != nil {
					log.Errorf("Cannot find the advertising campaign: %s", err)
					return err
				}

				for i := 0; i < 4; i++ {
					if i == 3 {
						title := fmt.Sprintf("ERROR: %s is not going to have the bid updated because of an API error", r.Name)
						message := fmt.Sprintf("ERROR: %s has a bid change, it is now %v, but the request to change failed %v times", r.Name, rp.Bid, i)
						var users Users
						users.MessageAdscoopsUsers(title, message)
						return nil
					}
					if err = ac.UpdateBid(rp.Bid); err != nil {
						log.Errorf("Cannot update bid: %s", err)
					} else {
						break
					}
				}

				title := fmt.Sprintf("%s has a bid change", r.Name)
				message := fmt.Sprintf("%s has a bid change, it is now: $%v", r.Name, rp.Bid)

				var users Users
				if err := users.MessageAdscoopsUsers(title, message); err != nil {
					return err
				}
			}

		}

		if r.AdvertisingDailySpend > 0 {
			sendNotification := false
			var rp RedirectPerformance
			err := AdscoopsDB.Where("redirect_id = ?", r.ID).Find(&rp).Error

			if err != nil {
				rp.RedirectID = r.ID
				rp.Budget = r.AdvertisingDailySpend
				rp.RequiredImpressions = uint(redirCampaignLimit)
				err := AdscoopsDB.Save(&rp).Error

				if err != nil {
					log.Errorf("Could not save redirect performance: %s", err)
					return err
				}
				sendNotification = true
			}

			log.Println("")
			log.Printf("redir id: %v", r.ID)

			log.Infof("requiredimpressions: %v", rp.RequiredImpressions)
			log.Infof("redirCampaignLimit: %v", redirCampaignLimit)

			log.Infof("budget: %v", rp.Budget)
			log.Infof("advertisingdailyspend: %v", r.AdvertisingDailySpend)

			log.Infof("rp asap: %v", rp.ASAPSpending)
			log.Infof("r asap: %v", r.EnableAdvertisingASAPSpending)

			log.Println("")

			if (rp.RequiredImpressions != uint(redirCampaignLimit)) || (rp.Budget != r.AdvertisingDailySpend) || rp.ASAPSpending != r.EnableAdvertisingASAPSpending || sendNotification {
				log.Infof("Daily spend should change because budget or impressions has changed")
				rp.Budget = uint((redirCampaignLimit / float64(rp.RequiredImpressions)) * float64(r.AdvertisingDailySpend))

				if rp.Budget == 0 && len(campaigns) == 0 {
					log.Error("Don't update the budget if no campaigns are enabled")
					return nil
				}
				rp.RequiredImpressions = uint(redirCampaignLimit)
				rp.ASAPSpending = r.EnableAdvertisingASAPSpending
				err := AdscoopsDB.Save(&rp).Error
				if err != nil {
					log.Errorf("Cannot update redirect performance in DB: %s", err)
					return err
				}

				AdscoopsDB.Model(&r).Update("advertising_daily_spend", rp.Budget)
				if err != nil {
					log.Errorf("Cannot update redirect budget in DB: %s", err)
					return err
				}

				title := fmt.Sprintf("%s has a budget change", r.Name)
				message := fmt.Sprintf("%s has a budget change, it is now: $%v with ASAP spending state set to: %v", r.Name, rp.Budget, rp.ASAPSpending)

				var users Users
				if err := users.MessageAdscoopsUsers(title, message); err != nil {
					log.Errorf("Cannot send message to adscoops users: %s", err)
					return err
				}

				if r.EnableAdvertisingSpendChange {
					var ac adonnetwork.Campaign
					if err = ac.Find(r.AdvertisingCampaignID); err != nil {
						log.Errorf("Cannot find the advertising campaign: %s", err)
						return err
					}

					budgetSpend := fmt.Sprintf("%v", rp.Budget)

					if rp.ASAPSpending {
						budgetSpend += ".99"
					}

					for i := 0; i < 4; i++ {
						if i == 3 {
							title := fmt.Sprintf("ERROR: %s has a budget change, but the API call failed", r.Name)
							message := fmt.Sprintf("%s has a budget change, it was attempted to be now: $%v with ASAP spending state set to: %v but the API call failed.", r.Name, rp.Budget, rp.ASAPSpending)
							var users Users
							users.MessageAdscoopsUsers(title, message)
							return nil
						}
						if err = ac.UpdateDailySpend(budgetSpend); err != nil {
							log.Errorf("Error updating daily spend budget: %s", err)
						} else {
							break
						}
					}

				}

			}
		}

	}

	return nil
}

func (r Redirect) UnpauseLinksInCampaigns() {
	AdscoopsDB.Exec(`UPDATE adscoop_urls
	JOIN adscoop_campaigns ON adscoop_urls.campaign_id = adscoop_campaigns.id
	SET adscoop_urls.deleted_at = NULL
	WHERE url LIKE '%/r/` + r.Hash + `%'
	AND adscoop_campaigns.enable_auto_redir_link_pause = 1;`)
}

func (r Redirect) PauseLinksInCampaigns() {
	AdscoopsDB.Exec(`UPDATE adscoop_urls
	JOIN adscoop_campaigns ON adscoop_urls.campaign_id = adscoop_campaigns.id
	SET adscoop_urls.deleted_at = NOW()
	WHERE url LIKE '%/r/` + r.Hash + `%'
	AND adscoop_campaigns.enable_auto_redir_link_pause = 1;`)
}

type Redirects []Redirect

func (r *Redirects) GetRecent() error {

	log.Infof("Getting recent redirects")
	now := time.Now().UTC()
	yesterday := now.Add(time.Duration(-24 * time.Hour))

	return AdscoopsDB.Select("adscoop_redirects.*").Table("adscoop_redirects").
		Joins("JOIN adscoop_trackings ON adscoop_trackings.redirect_id =  adscoop_redirects.id").
		Where("adscoop_trackings.timeslice >= ?", yesterday).
		Group("adscoop_trackings.redirect_id").Find(&r).Error
}

func (r *Redirects) TableName() string {
	return "adscoop_redirects"
}

func (r *Redirects) FindAll() error {
	return AdscoopsDB.Table("adscoop_redirects").Find(&r).Error
}

type RedirectCampaignRead struct {
	RedirectCampaign
	Name string
}

type RedirectCampaignReads []RedirectCampaignRead

func (r *RedirectCampaignReads) Find(id string) error {
	return AdscoopsDB.Table("adscoop_redirect_campaigns").Select("adscoop_campaigns.name, adscoop_redirect_campaigns.*").
		Joins("JOIN adscoop_campaigns ON adscoop_campaigns.ID = adscoop_redirect_campaigns.campaign_id").
		Where("adscoop_redirect_campaigns.redirect_id = ?", id).
		Find(&r).Error
}

func (r RedirectCampaignReads) Save(id string) error {

	err := AdscoopsDB.Unscoped().Where("redirect_id = ?", id).Delete(&RedirectCampaign{}).Error
	if err != nil {
		return err
	}

	for _, c := range r {
		var rc RedirectCampaign
		rc.CampaignID = c.CampaignID
		rc.RedirectID = c.RedirectID
		rc.Weight = c.Weight

		err := AdscoopsDB.Save(&rc).Error

		if err != nil {
			return err
		}

	}

	return nil
}

func (t *RedirectCampaignRead) TableName() string {
	return "adscoop_redirect_campaigns"
}

type RedirectCampaign struct {
	gorm.Model
	Weight     string `form:"weight"`
	CampaignID string `sql:"unique_index:redirect_id"`
	RedirectID string `sql:"unique_index:redirect_id"`
}

func (r *RedirectCampaign) TableName() string {
	return "adscoop_redirect_campaigns"
}

func (r Redirect) getCampaignsFromGroups(grabSoftPause bool) (campaigns []Campaign, err error) {
	now := getNow()

	var grabSoftPauseQuery string

	if grabSoftPause {
		grabSoftPauseQuery = "adscoop_campaigns.enable_soft_pause = 1 OR"
	}

	var sortMethod string
	if r.SortMethod == 0 {
		sortMethod = "adscoop_campaigns.campaign_group_weight DESC"
	} else {
		sortMethod = "adscoop_campaign_today_count.count / adscoop_campaigns.campaign_group_weight ASC"
	}

	err = AdscoopsDB.Select("adscoop_campaigns.*").Table("adscoop_campaigns").
		Joins(`JOIN adscoop_campaign_groups ON adscoop_campaign_groups.id = adscoop_campaigns.campaign_group_id
		JOIN adscoop_redirect_campaign_groups ON adscoop_redirect_campaign_groups.campaign_group_id = adscoop_campaign_groups.id
						 LEFT JOIN adscoop_campaign_today_count ON adscoop_campaign_today_count.campaign_id = adscoop_campaigns.id`).
		Where(fmt.Sprintf(`
				(adscoop_campaigns.deleted_at IS NULL or adscoop_campaigns.deleted_at <= '0001-01-02')
				AND ((%s (
					(adscoop_campaigns.paused = 0)
					AND (adscoop_campaigns.enable_start_stop_times = 0 OR ((adscoop_campaigns.disable_start_time = 1 OR adscoop_campaigns.start_datetime <= ?)
						AND (adscoop_campaigns.disable_end_time = 1 OR adscoop_campaigns.end_datetime > ?))
					)
				)) AND adscoop_redirect_campaign_groups.redirect_id = ?)
				`, grabSoftPauseQuery), now.String(), now.String(), r.ID).Group("adscoop_campaigns.id").Order(sortMethod).Find(&campaigns).Error

	if err != nil {
		log.Errorf("Error querying campaigns: %s", err)
		return
	}

	return
}

func (r Redirect) getCampaignsFromRedir(grabSoftPause bool) ([]Campaign, error) {
	now := getNow()

	var campaigns []Campaign
	var grabSoftPauseQuery string

	if grabSoftPause {
		grabSoftPauseQuery = "adscoop_campaigns.enable_soft_pause = 1 OR"
	}

	var sortMethod string
	if r.SortMethod == 0 {
		sortMethod = "adscoop_redirect_campaigns.weight DESC"
	} else {
		sortMethod = "adscoop_campaign_today_count.count / adscoop_redirect_campaigns.weight ASC"
	}

	err := AdscoopsDB.Select("adscoop_campaigns.*").Table("adscoop_campaigns").
		Joins(`JOIN adscoop_redirect_campaigns ON adscoop_redirect_campaigns.campaign_id = adscoop_campaigns.id
					 LEFT JOIN adscoop_campaign_today_count ON adscoop_campaign_today_count.campaign_id = adscoop_campaigns.id`).
		Where(fmt.Sprintf(`
			(adscoop_redirect_campaigns.deleted_at IS NULL or adscoop_redirect_campaigns.deleted_at <= '0001-01-02') AND
			(adscoop_campaigns.deleted_at IS NULL or adscoop_campaigns.deleted_at <= '0001-01-02')
			AND ((%s (
				(adscoop_campaigns.paused = 0)
				AND (adscoop_campaigns.enable_start_stop_times = 0 OR ((adscoop_campaigns.disable_start_time = 1 OR adscoop_campaigns.start_datetime <= ?)
					AND (adscoop_campaigns.disable_end_time = 1 OR adscoop_campaigns.end_datetime > ?))
				)
			)) AND adscoop_redirect_campaigns.redirect_id = ?)
			`, grabSoftPauseQuery), now.String(), now.String(), r.ID).Group("adscoop_campaigns.id").Order(sortMethod).Find(&campaigns).Error

	if err != nil {
		log.Errorf("Error querying campaigns: %s", err)
		return nil, err
	}

	return campaigns, err

}

func (r Redirect) GetLiveCampaigns(grabSoftPause bool) ([]Campaign, error) {
	var err error
	var campaigns []Campaign

	if r.CampaignSource == "0" {
		campaigns, err = r.getCampaignsFromRedir(grabSoftPause)
	}

	if r.CampaignSource == "1" {
		campaigns, err = r.getCampaignsFromGroups(grabSoftPause)
	}

	return campaigns, err
}

func (r Redirect) GetActiveCampaigns() ([]Campaign, error) {
	var activeCampaigns []Campaign
	campaigns, err := r.GetLiveCampaigns(true)
	if err != nil {
		return nil, err
	}
	for _, c := range campaigns {

		if !c.EnableSoftPause {

			if c.IsLimitReached() {
				continue
			}

			if !c.IsClientGoodStanding() {
				continue
			}
		}

		activeCampaigns = append(activeCampaigns, c)
	}

	return activeCampaigns, nil
}

type RedirectPerformance struct {
	gorm.Model
	RedirectID          uint
	Budget              uint
	RequiredImpressions uint
	ASAPSpending        bool `gorm:"column:asap_spending"`
	Bid                 string
}
