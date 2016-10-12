package adonnetwork

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

const (
	CampaignPauseState  = "pause"
	CampaignActiveState = "active"
)

type Campaign struct {
	gorm.Model
	AdvertiserID            int           `json:"advertiser_id"`
	BidType                 string        `json:"bid_type"`
	Bid                     string        `json:"bid"`
	CampaignMedia           []interface{} `json:"campaign_media" sql:"-"`
	Created                 string        `json:"created"`
	CurrentDaySpend         string        `json:"current_day_spend"`
	CurrentIntervalSpend    string        `json:"current_interval_spend"`
	CurrentMonthSpend       string        `json:"current_month_spend"`
	DailyCap                string        `json:"daily_cap"`
	EndDate                 string        `json:"end_date"`
	FeedType                string        `json:"feed_type"`
	FrequencyCap            int           `json:"frequency_cap"`
	FrequencyCount          int           `json:"frequency_count"`
	Geos                    []Geo         `json:"geos" sql:"-"`
	CampaignID              int           `json:"id" gorm:"column:campaign_id"`
	KeywordsPendingApproval int           `json:"keywords_pending_approval"`
	ListingsPendingApproval int           `json:"listings_pending_approval"`
	MinBid                  string        `json:"min_bid"`
	MonthlyCap              string        `json:"monthly_cap"`
	Name                    string        `json:"name"`
	StartDate               string        `json:"start_date"`
	Status                  string        `json:"status"`
	TotalCap                string        `json:"total_cap"`
	TotalSpend              string        `json:"total_spend"`
	Updated                 string        `json:"updated"`
}

type Campaigns []Campaign

func (c *Campaigns) FindAll() error {
	return AdscoopsDB.Table("adon_campaigns").Find(&c).Error
}

func (c Campaign) TableName() string {
	return "adon_campaigns"
}

func (c *Campaign) UpdateDailySpend(newSpend string) error {
	var jsonStr = []byte(fmt.Sprintf(`{"daily_cap": %s}`, newSpend))

	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiCall := makePutUrl(fmt.Sprintf(apiGetSingleCampaign, c.CampaignID))

	log.Printf("Going to call adon url: %s", apiCall)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) == "development" {
		return nil
	}

	req, err := http.NewRequest("PUT", apiCall, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Errorf("http new request error: %s", err)
		return err
	}

	_, err = client.Do(req)

	if err != nil {
		log.Errorf("client do error: %s", err)
		return err
	}

	return nil
}

func (c *Campaign) UpdateBid(bid string) error {
	var jsonStr = []byte(fmt.Sprintf(`{"bid": %s}`, bid))

	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiCall := makePutUrl(fmt.Sprintf(apiGetSingleCampaign, c.CampaignID))

	log.Printf("Going to call adon url: %s", apiCall)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) == "development" {
		return nil
	}

	req, err := http.NewRequest("PUT", apiCall, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Errorf("http new request error: %s", err)
		return err
	}

	_, err = client.Do(req)

	if err != nil {
		log.Errorf("client do error: %s", err)
		return err
	}

	return nil
}

func (c *Campaign) Pause() error {
	c.Status = CampaignPauseState
	if err := AdscoopsDB.Save(&c).Error; err != nil {
		return err
	}

	var jsonStr = []byte(fmt.Sprintf(`{"status": "%s"}`, CampaignPauseState))

	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiCall := makePutUrl(fmt.Sprintf(apiGetSingleCampaign, c.CampaignID))

	log.Printf("Going to call adon url: %s", apiCall)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) == "development" {
		return nil
	}

	req, err := http.NewRequest("PUT", apiCall, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Errorf("http new request error: %s", err)
		return err
	}

	_, err = client.Do(req)

	if err != nil {
		log.Errorf("client do error: %s", err)
		return err
	}

	return nil
}

func (c *Campaign) Find(id string) error {
	return AdscoopsDB.Where("campaign_id = ?", id).Find(&c).Error
}
func (c *Campaign) Resume() error {
	c.Status = CampaignActiveState
	if err := AdscoopsDB.Save(&c).Error; err != nil {
		return err
	}

	var jsonStr = []byte(fmt.Sprintf(`{"status": "%s"}`, CampaignActiveState))

	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiCall := makePutUrl(fmt.Sprintf(apiGetSingleCampaign, c.CampaignID))

	log.Printf("Going to call adon url: %s", apiCall)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) == "development" {
		return nil
	}

	req, err := http.NewRequest("PUT", apiCall, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Errorf("http new request error: %s", err)
		return err
	}

	res, err := client.Do(req)
	log.Infof("Status Code: %v", res.StatusCode)

	if err != nil {
		log.Errorf("client do error: %s", err)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Cannot read body: %s", err)
		return err
	}
	defer res.Body.Close()

	log.Printf("Body: %s", string(body))

	return nil
}

func UpdateAllCampaigns() error {
	var campaigns []Campaign
	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiCall := makeUrl(apiGetAllCampaigns)

	log.Printf("Going to call adon url: %s", apiCall)

	req, err := http.NewRequest("GET", apiCall, nil)

	if err != nil {
		log.Errorf("http new request error: %s", err)
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		log.Errorf("client do error: %s", err)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Errorf("read body error: %s", err)
		return err
	}

	err = json.Unmarshal(body, &campaigns)

	if err != nil {
		log.Errorf("unmarshal error, url:: %s, %s", err, apiCall)
		return err
	}

	tx := AdscoopsDB.Begin()

	err = tx.Exec("TRUNCATE adon_campaigns").Error

	if err != nil {
		errtx := tx.Rollback().Error
		if errtx != nil {
			log.Errorf("rollback error: %s", err)
			return errtx
		}
		log.Errorf("delete error: %s", err)
		return err
	}

	for _, c := range campaigns {
		err := tx.Save(&c).Error

		if err != nil {
			errtx := tx.Rollback().Error
			if errtx != nil {
				log.Errorf("rollback error: %s", err)
				return errtx
			}
			log.Errorf("delete error: %s", err)
			return err
		}
	}

	err = tx.Commit().Error

	if err != nil {
		log.Errorf("commit error: %s", err)
	}

	return err
}

type Geo struct {
	gorm.Model
	Abbreviation    string `json:"abbreviation"`
	AbbreviationAlt string `json:"abbreviation_alt"`
	ID              int    `json:"id" gorm:"column:geo_id"`
	Level           int    `json:"level"`
	Name            string `json:"name"`
	Number          int    `json:"number"`
}

func (c Geo) TableName() string {
	return "adon_geos"
}
