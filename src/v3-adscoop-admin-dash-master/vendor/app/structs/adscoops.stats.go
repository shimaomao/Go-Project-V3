package structs

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

)

func getDayStart() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	local := time.Now().In(loc)

	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}
func getHourStart() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	local := time.Now().In(loc)

	return time.Date(local.Year(), local.Month(), local.Day(), local.Hour(), 0, 0, 0, loc)
}

func getNow() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	return time.Now().In(loc)
}

/*GetDailyImpressionsCountByVertical gets the daily counts broken down by loads, engagements and clicks */
func (s *Trackings) GetDailyImpressionsCountByVertical(trackingMethod uint) {
	today := getDayStart()

	currentHour := getHourStart()

	var prevTrackings TrackingsRaw
	var hourTrackings TrackingsRaw

	prevImpsCountByVertKey := fmt.Sprintf("imps_count_by_vert_%v_%s", trackingMethod, currentHour)
	prevImpsCountByVert, found := Cache.Get(prevImpsCountByVertKey)

	if !found {
		impsrevByVertQuery(trackingMethod, today.UTC(), currentHour.UTC()).Find(&prevTrackings)

		Cache.Set(prevImpsCountByVertKey, prevTrackings, time.Hour)
	} else {
		prevTrackings = prevImpsCountByVert.(TrackingsRaw)
	}

	impsrevByVertQuery(trackingMethod, currentHour.UTC(), currentHour.Add(time.Hour).UTC()).Find(&hourTrackings)

	s.Engagement = prevTrackings.Engagement + hourTrackings.Engagement
	s.Count = prevTrackings.Count + hourTrackings.Count
	s.Cpc = prevTrackings.Cpc + hourTrackings.Cpc
}

/*GetDailyImpressionsCount gets the click count for the day */
func (s *Trackings) GetDailyImpressionsCount() {
	today := getDayStart()

	var prevTrackings TrackingsRaw
	var hourTrackings TrackingsRaw

	currentHour := getHourStart()

	prevImpsCountKey := fmt.Sprintf("imps_rev_count_%s", currentHour)

	prevImpsCount, found := Cache.Get(prevImpsCountKey)

	if !found {
		impsrevQuery(today.UTC(), currentHour.UTC()).Find(&prevTrackings)

		Cache.Set(prevImpsCountKey, prevTrackings, time.Hour)
	} else {
		prevTrackings = prevImpsCount.(TrackingsRaw)
	}

	impsrevQuery(currentHour.UTC(), currentHour.Add(time.Hour).UTC()).Find(&hourTrackings)

	s.Engagement = prevTrackings.Engagement + hourTrackings.Engagement
	s.Count = prevTrackings.Count + hourTrackings.Count
	s.Cpc = prevTrackings.Cpc + hourTrackings.Cpc
}

type TrackingRows []Trackings

func (s *TrackingRows) GetDailyImpressionsCountByVertical(trackingMethod uint) {
	var todayStats []Trackings
	var hourStats []Trackings

	today := getDayStart()
	currentHour := getHourStart()

	prevImpsGroupByVertKey := fmt.Sprintf("imps_rev_count_group_%v_%s", trackingMethod, currentHour)

	prevImpsGroupByVert, found := Cache.Get(prevImpsGroupByVertKey)

	if !found {
		impsrevGroupByVertQuery(trackingMethod, today.UTC(), currentHour.UTC()).Find(&todayStats)

		Cache.Set(prevImpsGroupByVertKey, todayStats, time.Hour)
	} else {
		todayStats = prevImpsGroupByVert.([]Trackings)
	}

	impsrevGroupByVertQuery(trackingMethod, currentHour.UTC(), currentHour.Add(time.Hour).UTC()).Find(&hourStats)

	for _, t := range todayStats {
		*s = append(*s, t)
	}

	for _, t := range hourStats {
		*s = append(*s, t)
	}

}

func (t *TrackingRows) ForClient(clientID uint, trackingMethod uint) (uint, error) {
	var tr []Trackings
	var trDay []Trackings
	today := getDayStart()

	currentHour := getHourStart()

	prevResultsKey := fmt.Sprintf("client_stats_%v_%v_%s", trackingMethod, clientID, currentHour.String())

	prevResults, found := Cache.Get(prevResultsKey)

	if !found {
		query := statsQuery(today.UTC(), trackingMethod, clientID)

		err := query.Find(&trDay).Error

		if err != nil {
			return 0, err
		}

		Cache.Set(prevResultsKey, trDay, time.Hour)
	} else {
		trDay = prevResults.([]Trackings)
	}

	query := statsQuery(currentHour.UTC(), trackingMethod, clientID)
	var impCount uint

	err := query.Find(&tr).Error

	for _, r := range trDay {
		if r.Timeslice.UTC().Hour() == currentHour.UTC().Hour() {
			continue
		}
		impCount += r.Load

		*t = append(*t, r)
	}

	for _, r := range tr {
		impCount += r.Load

		*t = append(*t, r)
	}

	return impCount, err
}

type Trackings struct {
	Timeslice  time.Time
	RedirectID uint
	UrlID      uint
	Count      uint
	Cpc        float64
	Engagement uint
	Load       uint
	TimeOnSite string
}

type TrackingsRaw struct {
	Engagement uint
	Count      uint
	Cpc        float64
}

func (t Trackings) TableName() string {
	return "adscoop_trackings"
}

func statsQuery(timeslice time.Time, trackingMethod uint, clientID uint) *gorm.DB {
	query := AdscoopsDB.
		Select("DISTINCT C2.timeslice, ROUND(SUM(C2.daily_limit)) as engagement, ROUND(SUM(C2.count)) as count, SUM(C2.count) as `load`, ROUND(SUM(C2.revenue)) as cpc").
		Table("adscoop_campaigns").
		Joins(fmt.Sprintf(`
	LEFT OUTER JOIN(
			SELECT DISTINCT timeslice, adscoop_campaigns.id as campaign_id, adscoop_campaigns.daily_imps_limit as daily_limit, SUM(CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
							 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement_count
							 WHEN adscoop_campaigns.tracking_method = 2 THEN temp_stats.load_count
						END) as count, SUM(temp_stats.revenue) as revenue
			FROM temp_stats
			JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
			JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
			WHERE timeslice >= '%s'
			GROUP BY adscoop_campaigns.id, timeslice
		) C2 on adscoop_campaigns.id = C2.campaign_id
	`, timeslice)).
		Where("tracking_method = ? AND count >= 1 AND adscoop_campaigns.client_id = ?", trackingMethod, clientID).Group("DATE(C2.timeslice), HOUR(C2.timeslice)")

	return query
}

func impsrevQuery(timeslice time.Time, endTimeslice time.Time) *gorm.DB {

	return AdscoopsDB.LogMode(true).Select("DISTINCT SUM(C2.daily_limit) as engagement, SUM(C2.count) as count, SUM(C2.revenue) as cpc").Table("adscoop_campaigns").Joins(fmt.Sprintf(`
LEFT OUTER JOIN(
		SELECT DISTINCT adscoop_campaigns.id as campaign_id, adscoop_campaigns.daily_imps_limit as daily_limit, SUM(CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
						 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement_count
						 WHEN adscoop_campaigns.tracking_method = 2 THEN temp_stats.load_count
					END) as count, SUM(temp_stats.revenue) as revenue
		FROM temp_stats
		JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
		JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
		WHERE timeslice >= '%s' AND timeslice < '%s'
		GROUP BY adscoop_campaigns.id
	) C2 on adscoop_campaigns.id = C2.campaign_id
`, timeslice.String(), endTimeslice.String()))
}

func impsrevByVertQuery(trackingMethod uint, timeslice time.Time, endTimeslice time.Time) *gorm.DB {
	return AdscoopsDB.Select("DISTINCT SUM(C2.daily_limit) as engagement, SUM(C2.count) as count, SUM(C2.revenue) as cpc").Table("adscoop_campaigns").Joins(fmt.Sprintf(`
	LEFT OUTER JOIN(
			SELECT DISTINCT adscoop_campaigns.id as campaign_id, adscoop_campaigns.daily_imps_limit as daily_limit, SUM(CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
							 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement_count
							 WHEN adscoop_campaigns.tracking_method = 2 THEN temp_stats.load_count
						END) as count, SUM(temp_stats.revenue) as revenue
			FROM temp_stats
			JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
			JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
			WHERE timeslice >= '%s' AND timeslice < '%s'
			GROUP BY adscoop_campaigns.id
		) C2 on adscoop_campaigns.id = C2.campaign_id
	`, timeslice.String(), endTimeslice.String())).
		Where("tracking_method = ?", trackingMethod)
}

func impsrevGroupByVertQuery(trackingMethod uint, timeslice time.Time, endTimeslice time.Time) *gorm.DB {
	return AdscoopsDB.
		Select("DISTINCT C2.timeslice, SUM(C2.daily_limit) as engagement, SUM(C2.count) as count, SUM(C2.revenue) as cpc").
		Table("adscoop_campaigns").
		Joins(fmt.Sprintf(`
	LEFT OUTER JOIN(
			SELECT DISTINCT timeslice, adscoop_campaigns.id as campaign_id, adscoop_campaigns.daily_imps_limit as daily_limit, SUM(CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
							 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement_count
							 WHEN adscoop_campaigns.tracking_method = 2 THEN temp_stats.load_count
						END) as count, SUM(temp_stats.revenue) as revenue
			FROM temp_stats
			JOIN adscoop_urls ON adscoop_urls.id = temp_stats.url_id
			JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
			WHERE timeslice >= '%s' AND timeslice < '%s'
			GROUP BY adscoop_campaigns.id, timeslice
		) C2 on adscoop_campaigns.id = C2.campaign_id
	`, timeslice.String(), endTimeslice.String())).
		Where("tracking_method = ? AND count >= 1", trackingMethod).
		Group("DATE(C2.timeslice), HOUR(C2.timeslice)")
}
