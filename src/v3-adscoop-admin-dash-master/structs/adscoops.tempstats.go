package structs

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// TempStats is used for stats saved in the past 30 hours
type TempStats struct {
	Timeslice       time.Time
	RedirectID      uint
	Revenue         float64
	Count           uint
	EngagementCount uint
	LoadCount       uint
	UrlID           uint
}

type TempStatsRead struct {
	TempStats
	AllCount uint
}

type TempClientStats struct {
	TempStatsRead
	ClientName string
	ChartColor string
	ClientID   uint
}

type MultiTempStats []TempStatsRead
type MultiClientTempStats []TempClientStats

func (m *MultiClientTempStats) Today() {
	now := time.Now().UTC()
	hoursAgo := now.Add(time.Duration(-lookbackInMinutes * time.Minute))

	AdscoopsDB.Select(`adscoop_clients.id as client_id,
						 adscoop_clients.name as client_name,
						 adscoop_clients.chart_color as chart_color,
						 timeslice,
					   SUM(revenue) as revenue,
						 SUM(count + engagement_count + load_count) AS all_count,
						 SUM(count) as count,
						 SUM(engagement_count) as engagement_count,
						 SUM(load_count) as load_count`).
		Table("temp_stats").
		Joins(`JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
	JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
	JOIN adscoop_clients ON adscoop_clients.id = adscoop_campaigns.client_id`).
		Where("timeslice < ? AND timeslice >= ?", now.String(), hoursAgo.String()).
		Group("adscoop_clients.id, timeslice").Find(&m)
}

func (m *MultiTempStats) Yesterday() {
	now := time.Now().UTC()
	yesterday := now.Add(time.Duration(-24 * time.Hour))

	hoursAgo := yesterday.Add(time.Duration(-30 * time.Minute))
	AdscoopsRealtimeDB.Select("timeslice, SUM(revenue) as revenue, SUM(count + engagement_count + load_count) AS all_count, SUM(count) as count, SUM(engagement_count) as engagement_count, SUM(load_count) as load_count").Table("temp_stats").Where("timeslice < ? AND timeslice >= ?", yesterday.String(), hoursAgo.String()).Group("timeslice").Order("timeslice DESC").Limit(30).Find(&m)
}

func (m *MultiTempStats) Today() {
	now := time.Now().UTC()

	hoursAgo := now.Add(time.Duration(-30 * time.Minute))
	AdscoopsRealtimeDB.Select("timeslice, SUM(revenue) as revenue, SUM(count + engagement_count + load_count) AS all_count, SUM(count) as count, SUM(engagement_count) as engagement_count, SUM(load_count) as load_count").Table("temp_stats").Where("timeslice < ? AND timeslice >= ?", now.String(), hoursAgo.String()).Group("timeslice").Order("timeslice DESC").Limit(30).Find(&m)
}

type TempStatsContainer struct {
	lock sync.RWMutex
	list map[string]TempStats
}

func (t *TempStatsContainer) Add(ts *TempStats) {
	var tskey = fmt.Sprintf("ts%v%v", ts.RedirectID, ts.UrlID)

	t.lock.Lock()

	var nts = t.list[tskey]
	if nts.RedirectID == 0 {
		nts.RedirectID = ts.RedirectID
		nts.UrlID = ts.UrlID
	}
	nts.Count += ts.Count
	nts.EngagementCount += ts.EngagementCount
	nts.LoadCount += ts.LoadCount
	nts.Revenue += ts.Revenue

	t.list[tskey] = nts

	log.Printf("reading: %+v", ts)
	log.Printf("adding: %+v", nts)

	t.lock.Unlock()
}

func (t *TempStatsContainer) Save() {
	t.list = make(map[string]TempStats)
	ticker := time.Tick(time.Duration(5) * time.Second)

	for _ = range ticker {
		t.Push()
	}
}

func (t *TempStatsContainer) Push() {
	t.lock.Lock()
	list := t.list
	t.list = make(map[string]TempStats)
	t.lock.Unlock()

	for _, u := range list {
		timeslice := time.Date(time.Now().UTC().Year(),
			time.Now().UTC().Month(),
			time.Now().UTC().Day(),
			time.Now().UTC().Hour(),
			time.Now().UTC().Minute(), 0, 0, time.UTC)

		_, err := AdscoopsRealtimeDB.DB().Exec(`INSERT temp_stats(timeslice, redirect_id, url_id, revenue, count, engagement_count, load_count)
			VALUES(?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				count = count + ?,
				engagement_count = engagement_count + ?,
				load_count = load_count + ?,
				revenue = revenue + ?`,
			timeslice, u.RedirectID, u.UrlID, u.Revenue, u.Count, u.EngagementCount, u.LoadCount,
			u.Count, u.EngagementCount, u.LoadCount, u.Revenue)

		if err != nil {
			log.Println("stats err", err)
		}
	}
}
