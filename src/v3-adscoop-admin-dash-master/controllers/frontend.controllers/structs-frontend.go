package frontendControllers

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	//"strconv"
	"app/adscoops.caches"
	"app/structs"
	"sync"
	"time"

	"strconv"
)

const TRACKING_DATA_KEY = "td%v%v%s%s"

var tsc = structs.TempStatsContainer{}

type TomlConfig struct {
	MysqlConfig string
}

type ErrorXml struct {
	XMLName xml.Name `xml:"error" json:"-"`
	Message string   `xml:"message"`
}

type ReturnXml struct {
	XMLName xml.Name         `xml:"results" json:"-"`
	Link    []*ReturnLinkXml `xml:"link" json:"link"`
}

type ReturnLinkXml struct {
	Title string `xml:"title" json:"title"`
	Link  string `xml:"link" json:"link"`
	Cpc   string `xml:"cpc" json:"cpc"`
}

type ClickStats struct {
	lock sync.RWMutex
	list map[string]trackingData
}

type UrlTrackingMethod struct {
	lock sync.RWMutex
	list map[uint]string
}

var urlTrackingMethod UrlTrackingMethod

type trackingData struct {
	RedirectID       uint
	UrlID            uint
	CPC              string
	UniqueIdentifier string
	Count            int64
	EngagementCount  int64
	LoadCount        int64
	TimeOnSite       float64
	TimeOnSiteCount  float64
}

func (t *ClickStats) Add(ast *AdscoopTracking, trackingType int64) {
	c := adscoopsCaches.RedisPool.Get()
	defer c.Close()

	var tskey = fmt.Sprintf(TRACKING_DATA_KEY, ast.RedirectId, ast.UrlId, ast.Cpc, ast.UniqueIdentifier)

	t.lock.Lock()

	var ts structs.TempStats
	ts.RedirectID = uint(ast.RedirectId)

	urlTrackingMethod.lock.RLock()
	trackingMethod := urlTrackingMethod.list[ast.UrlId]
	urlTrackingMethod.lock.RUnlock()

	if trackingMethod == fmt.Sprintf("%v", trackingType) {
		if s, err := strconv.ParseFloat(ast.Cpc, 64); err == nil {
			ts.Revenue = s
		}
	}

	ts.UrlID = uint(ast.UrlId)

	if t.list[tskey].RedirectID == 0 {
		var td trackingData
		td.UrlID = ast.UrlId
		td.RedirectID = ast.RedirectId
		td.CPC = ast.Cpc
		td.UniqueIdentifier = ast.UniqueIdentifier

		if trackingType == 0 {
			ts.Count = 1
			td.Count = 1
		}
		if trackingType == 1 {
			ts.EngagementCount = 1
			td.EngagementCount = 1
		}
		if trackingType == 2 {
			ts.LoadCount = 1
			td.LoadCount = 1
		}
		if trackingType == 3 {
			td.TimeOnSiteCount = 1
			td.TimeOnSite = ast.TimeOnSite
		}

		t.list[tskey] = td
	} else {
		var td = t.list[tskey]
		if trackingType == 0 {
			ts.Count = 1
			td.Count += 1
		}
		if trackingType == 1 {
			ts.EngagementCount = 1
			td.EngagementCount += 1
		}
		if trackingType == 2 {
			ts.LoadCount = 1
			td.LoadCount += 1
		}
		if trackingType == 3 {
			td.TimeOnSiteCount += 1
			td.TimeOnSite += ast.TimeOnSite
		}
		t.list[tskey] = td
	}

	if trackingType != 3 {
		tsc.Add(&ts)
		//Add(&ts)
	}

	t.lock.Unlock()
}

func (t *ClickStats) Save() {
	ticker := time.Tick(time.Duration(60) * time.Second)

	for tc := range ticker {
		t.Push()

		log.Printf("PUshed at: %v\n", tc)
	}
}

func (t *ClickStats) Push() {
	t.lock.Lock()
	list := t.list
	t.list = make(map[string]trackingData)
	t.lock.Unlock()

	for x, u := range list {
		log.Println("saving for ", x)
		log.Println("adding new count ", u.Count)
		log.Println("adding new engagements ", u.EngagementCount)
		log.Println("adding new load ", u.LoadCount)

		minutes := float64(time.Now().UTC().Minute()) / 30
		minRound := int(math.Floor(minutes)) * 30

		timeslice := time.Date(time.Now().UTC().Year(),
			time.Now().UTC().Month(),
			time.Now().UTC().Day(),
			time.Now().UTC().Hour(),
			minRound, 0, 0, time.UTC)

		_, err := db.DB().Exec(`INSERT adscoop_trackings(timeslice, redirect_id, url_id, cpc, unique_identifier, count, engagement, adscoop_trackings.load, time_on_site, time_on_site_count)
			VALUES( ? , ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				count = count + ?,
				engagement = engagement + ?,
				adscoop_trackings.load = adscoop_trackings.load + ?,
				time_on_site = time_on_site + ?,
				time_on_site_count = time_on_site_count + ?`, timeslice, u.RedirectID, u.UrlID, u.CPC, u.UniqueIdentifier,
			u.Count, u.EngagementCount, u.LoadCount, u.TimeOnSite, u.TimeOnSiteCount, u.Count, u.EngagementCount, u.LoadCount, u.TimeOnSite, u.TimeOnSiteCount)

		if err != nil {
			log.Println("stats err", err)
		}
	}

	log.Println("Stats have been saved to DB")
}

type AdscoopRedirect struct {
	Id                      int64
	Hash                    string
	Iframe                  int64
	Min                     int64
	Max                     int64
	AutoRefresh             int64
	StripReferrer           int64
	RedirType               int64
	Name                    string
	StripQueryString        int64
	BapiScoring             int64
	ForceRefresh            bool `form:"force_refresh" json:"force_refresh"`
	SortMethod              uint
	LockWhitelistId         int64
	LockWhitelistUrls       []AdscoopWhitelistUrl
	LockUseragentId         int64
	ForceHost               int64
	LockWhitelistReverse    uint   `form:"lock_whitelist_reverse" json:"lock_whitelist_reverse"`
	LockUseragentReverse    uint   `form:"lock_useragent_reverse" json:"lock_useragent_reverse"`
	ForceHostString         string `sql:"-"`
	BbsiPath                string
	LockUseragents          []AdscoopWhitelistUseragent
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               time.Time
	ScoringTimeout          uint   `form:"scoring_timeout" json:"scoring_timeout"`
	ScoringRedirectEnabled  bool   `form:"scoring_redirect_enabled" json:"scoring_redirect_enabled"`
	ScoringRedirectOverride string `form:"scoring_redirect_override" json:"scoring_redirect_override"`
}

type AdscoopFeed struct {
	Id        int64
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Redirects []AdscoopRedirect
}

type AdscoopCampaign struct {
	Id                   int64
	Cpc                  string
	DailyImpsLimit       int64
	Urls                 []*AdscoopCampaignUrl
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	Type                 int64
	XmlType              int64
	TrackingMethod       int64
	EnableUnloadTracking bool
}

type AdscoopCampaignUrl struct {
	Id         int64
	CampaignId int64
	Weight     int64
	Url        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

func (a AdscoopCampaignUrl) TableName() string {
	return "adscoop_urls"
}

type AdscoopRedirectQuerystring struct {
	Id             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	RedirectId     int64
	QueryStringKey string
}

type RetData struct {
	AdscoopRedirect structs.Redirect
	BapiScoring     string
	QueryString     string
	AllowRefresh    bool
}

type AdscoopTracking struct {
	Timeslice        time.Time
	RedirectId       uint
	UrlId            uint
	UniqueIdentifier string
	Cpc              string
	TimeOnSite       float64
	TimeOnSiteCount  float64
}

func (a *AdscoopTracking) Track() {
	a.Timeslice = time.Date(time.Now().UTC().Year(),
		time.Now().UTC().Month(),
		time.Now().UTC().Day(),
		time.Now().UTC().Hour(),
		0, 0, 0, time.UTC)

	db.DB().Exec(`INSERT adscoop_trackings(timeslice, redirect_id, url_id, cpc, count, engagement, adscoop_trackings.load)
		VALUES( ? , ?, ?, ?, 1, 0, 0)
		ON DUPLICATE KEY UPDATE
			count = count + 1`, a.Timeslice, a.RedirectId, a.UrlId, a.Cpc)
}

func (a *AdscoopTracking) TrackEngagement() {
	a.Timeslice = time.Date(time.Now().UTC().Year(),
		time.Now().UTC().Month(),
		time.Now().UTC().Day(),
		time.Now().UTC().Hour(),
		0, 0, 0, time.UTC)

	_, err := db.DB().Exec(`INSERT adscoop_trackings(timeslice, redirect_id, url_id, cpc, count, engagement, adscoop_trackings.load)
		VALUES( ? , ?, ?, ?, 0, 1, 0)
		ON DUPLICATE KEY UPDATE
			engagement = engagement + 1`, a.Timeslice, a.RedirectId, a.UrlId, a.Cpc)

	if err != nil {
		log.Println("err", err)
	}
}

func (a *AdscoopTracking) TrackLoad() {
	a.Timeslice = time.Date(time.Now().UTC().Year(),
		time.Now().UTC().Month(),
		time.Now().UTC().Day(),
		time.Now().UTC().Hour(),
		0, 0, 0, time.UTC)

	_, err := db.DB().Exec(`INSERT adscoop_trackings(timeslice, redirect_id, url_id, cpc, count, engagement, adscoop_trackings.load)
		VALUES( ? , ?, ?, ?, 0, 0, 1)
		ON DUPLICATE KEY UPDATE
			adscoop_trackings.load = adscoop_trackings.load + 1`, a.Timeslice, a.RedirectId, a.UrlId, a.Cpc)

	if err != nil {
		log.Println("err", err)
	}
}

type AdscoopWhitelist struct {
	Id   int64    `form:"id"`
	Name string   `form:"name"`
	Urls []string `form:"url[]" sql:"-"`
}

type AdscoopWhitelistUrl struct {
	Id                 int64  `form:"id"`
	Url                string `form:"url"`
	AdscoopWhitelistId int64  `form:"whitelist_id"`
}

type AdscoopWhitelistUseragentGroup struct {
	Id         int64    `form:"id"`
	Name       string   `form:"name"`
	Useragents []string `form:"ua[]" sql:"-"`
}

type AdscoopWhitelistUseragent struct {
	Id                               int64  `form:"id"`
	Useragent                        string `form:"ua"`
	AdscoopWhitelistUseragentGroupId int64  `form:"whitelist_useragent_group"`
}

type AdscoopHost struct {
	Id   int64
	Host string
}
