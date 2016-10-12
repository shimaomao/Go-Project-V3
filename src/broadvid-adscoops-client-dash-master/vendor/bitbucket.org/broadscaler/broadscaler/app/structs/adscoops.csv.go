package structs

import "time"

type AdscoopClientCsvPost struct {
	ClientId  int64  `form:"client_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type AdscoopCsvStats struct {
	CampaignName               string
	Title                      string
	Url                        string
	Impressions                int
	Timeslice                  time.Time
	Verifieds                  int
	Loads                      int
	TrackingMethod             int64
	ImpressionsString          string
	Cost                       float64
	CostVerified               float64
	CostLoad                   float64
	CostString                 string
	Cpc                        string
	UniqueIdentifier           string
	EnableReportAccountBalance bool
	ClientID                   uint
}
