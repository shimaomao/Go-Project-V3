package structs

import (
	"encoding/xml"
	//"fmt"
	"html/template"
	//"os"
	"time"

	"github.com/jinzhu/gorm"
)

type AdscoopClientEmail struct {
	gorm.Model
	Email    string `sql:"unique_index:client_id"`
	ClientID uint   `sql:"unique_index:client_id"`
}

type AdscoopClientCampaignEmail struct {
	gorm.Model
	Email    string `sql:"unique_index:client_id"`
	ClientID uint   `sql:"unique_index:client_id"`
}

type AdscoopClientTransaction struct {
	gorm.Model
	ClientID      uint   `form:"client_id"`
	AmountCharged string `sq:"type:decimal(11,2)" form:"amount_charged"`
	TransactionId string
	Successful    bool
	Attempts      uint
}

type AdscoopClientUser struct {
	gorm.Model
	Name         string `form:"name"`
	Email        string `form:"email"`
	Password     string `form:"password"`
	ClientID     uint   `form:"client_id"`
	UserPolicyID uint   `form:"user_policy_id"`
}

type AdscoopClient struct {
	gorm.Model
	Name                        string `form:"name"`
	Paused                      bool   `form:"paused"`
	DefaultCpc                  string `sql:"type:decimal(11,4)" form:"default_cpc"`
	DailyImpsLimit              uint   `form:"daily_imps_limit"`
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
	Emails                      []string  `form:"email[]" json:"emails" sql:"-"`
	CampaignEmails              []string  `form:"campaign_email[]" json:"campaignEmails" sql:"-"`
	Transactions                []AdscoopClientTransaction
	EnableClientLogin           bool `form:"enable_client_login"`
	ClientSchedulesPendApproval bool `form:"client_schedules_pend_approval"`
	ShowMtdSpendInReport        bool `form:"show_mtd_spend_in_report"`
	InGoodStanding              bool `form:"in_good_standing"`
}

type AdscoopUserPolicy struct {
	gorm.Model
	Name                                          string `form:"name"`
	NameHidden                                    bool   `form:"name_hidden"`
	NameReadOnly                                  bool   `form:"name_readonly"`
	SourceHidden                                  bool   `form:"source_hidden"`
	SourceReadOnly                                bool   `form:"source_readonly"`
	DailyBudgetHidden                             bool   `form:"daily_budget_hidden"`
	DailyBudgetReadOnly                           bool   `form:"daily_budget_readonly"`
	CPCHidden                                     bool   `form:"cpc_hidden"`
	CPCReadOnly                                   bool   `form:"cpc_readonly"`
	WeightVarianceHidden                          bool   `form:"weight_variance_hidden"`
	WeightVarianceReadOnly                        bool   `form:"weight_variance_readonly"`
	TrackingMethodHidden                          bool   `form:"tracking_method_hidden"`
	TrackingMethodReadOnly                        bool   `form:"tracking_method_readonly"`
	PausedHidden                                  bool   `form:"paused_hidden"`
	PausedReadOnly                                bool   `form:"paused_readonly"`
	EnableMonCheckingHidden                       bool   `form:"enable_mon_checking_hidden"`
	EnableMonCheckingReadOnly                     bool   `form:"enable_mon_checking_readonly"`
	EnablePerformanceBasedPauseHidden             bool   `form:"enable_performance_based_pause_hidden"`
	EnablePerformanceBasedPauseReadOnly           bool   `form:"enable_performance_based_pause_readonly"`
	EnablePerformanceBasedPauseNotifyOnlyHidden   bool   `form:"enable_performance_based_pause_notify_only_hidden"`
	EnablePerformanceBasedPauseNotifyOnlyReadOnly bool   `form:"enable_performance_based_pause_notify_only_readonly"`
	EPBPCompareAHidden                            bool   `form:"epbp_compare_a_hidden"`
	EPBPCompareAReadOnly                          bool   `form:"epbp_compare_a_readonly"`
	EPBPLessThanHidden                            bool   `form:"epbp_less_than_hidden"`
	EPBPLessThanReadOnly                          bool   `form:"epbp_less_than_readonly"`
	EPBPCompareBHidden                            bool   `form:"epbp_compare_b_hidden"`
	EPBPCompareBReadOnly                          bool   `form:"epbp_compare_b_readonly"`
	FlightStartTimeHidden                         bool   `form:"flight_start_time_hidden"`
	FlightStartTimeReadOnly                       bool   `form:"flight_start_time_readonly"`
	DisableOnlyStartTimeHidden                    bool   `form:"disable_only_start_time_hidden"`
	DisableOnlyStartTimeReadOnly                  bool   `form:"disable_only_start_time_readonly"`
	FlightEndTimeHidden                           bool   `form:"flight_end_time_hidden"`
	FlightEndTimeReadOnly                         bool   `form:"flight_end_time_readonly"`
	DisableOnlyEndTimeHidden                      bool   `form:"disable_only_end_time_hidden"`
	DisableOnlyEndTimeReadOnly                    bool   `form:"disable_only_end_time_readonly"`
	EnableStartStopTimesHidden                    bool   `form:"enable_start_stop_times_hidden"`
	EnableStartStopTimesReadOnly                  bool   `form:"enable_start_stop_times_readonly"`
	UrlsHidden                                    bool   `form:"urls_hidden"`
	UrlsReadOnly                                  bool   `form:"urls_readonly"`
	MacrosHidden                                  bool   `form:"macros_hidden"`
	AddCopyCampaignEnabled                        bool   `form:"addcopy_campaign_enabled"`
	HideUrlUpdates                                bool   `form:"hide_url_updates"`
	EnableBusinessHours                           bool   `form:"enable_business_hours"`
	MondayStartHour                               int    `form:"monday_start_hour"`
	MondayEndHour                                 int    `form:"monday_end_hour"`
	MondayBlackout                                bool   `form:"monday_blackout"`
	TuesdayStartHour                              int    `form:"tuesday_start_hour"`
	TuesdayEndHour                                int    `form:"tuesday_end_hour"`
	TuesdayBlackout                               bool   `form:"tuesday_blackout"`
	WednesdayStartHour                            int    `form:"wednesday_start_hour"`
	WednesdayEndHour                              int    `form:"wednesday_end_hour"`
	WednesdayBlackout                             bool   `form:"wednesday_blackout"`
	ThursdayStartHour                             int    `form:"thursday_start_hour"`
	ThursdayEndHour                               int    `form:"thursday_end_hour"`
	ThursdayBlackout                              bool   `form:"thursday_blackout"`
	FridayStartHour                               int    `form:"friday_start_hour"`
	FridayEndHour                                 int    `form:"friday_end_hour"`
	FridayBlackout                                bool   `form:"friday_blackout"`
	SaturdayStartHour                             int    `form:"saturday_start_hour"`
	SaturdayEndHour                               int    `form:"saturday_end_hour"`
	SaturdayBlackout                              bool   `form:"saturday_blackout"`
	SundayStartHour                               int    `form:"sunday_start_hour"`
	SundayEndHour                                 int    `form:"sunday_end_hour"`
	SundayBlackout                                bool   `form:"sunday_blackout"`
	AutoApproveDelay                              int    `form:"auto_approve_delay"`
	AppendRcHidden                                bool   `form:"append_rc_hidden"`
	AppendRcReadOnly                              bool   `form:"append_rc_readonly"`
}

type AdscoopFeedRedirect struct {
	gorm.Model
	Weight     uint   `form:"weight"`
	FeedId     uint   `form:"feed_id"`
	MinCpc     string `sql:"type:decimal(11,4)" form:"min_cpc"`
	RedirectID uint   `form:"redirect_id"`
}

type AdscoopFeed struct {
	gorm.Model
	Name      string `form:"name" json:"name"`
	Paused    bool   `form:"paused"`
	Hash      string
	Redirects []AdscoopFeedRedirect
}

type AdscoopHost struct {
	gorm.Model
	Host string
}

type AdscoopPaymentHash struct {
	gorm.Model
	Hash     string
	ClientID uint
}

type AdscoopRedirectCampaign struct {
	gorm.Model
	Weight     uint `form:"weight"`
	CampaignID uint `form:"campaign_id" sql:"unique_index:redirect_id"`
	RedirectID uint `form:"redirect_id" sql:"unique_index:redirect_id"`
}

type AdscoopRedirectCampaignRead struct {
	AdscoopRedirectCampaign
	Name string
}

func (t *AdscoopRedirectCampaignRead) TableName() string {
	return "adscoop_redirect_campaigns"
}

type AdscoopRedirectQuerystring struct {
	gorm.Model
	RedirectID     uint   `sql:"unique_index:redirect_id"`
	QueryStringKey string `sql:"unique_index:redirect_id"`
}

type AdscoopRedirect struct {
	gorm.Model
	Name                    string   `form:"name" json:"name"`
	Hash                    string   `sql:"unique_index"`
	Min                     uint     `form:"min" json:"min"`
	Max                     uint     `form:"max" json:"max"`
	Iframe                  uint     `form:"iframe" json:"iframe"`
	RedirType               uint     `form:"redir_type" json:"redir_type"`
	BapiScoring             uint     `form:"bapi_scoring" json:"bapi_scoring"`
	LockWhitelistId         uint     `form:"lock_whitelist_id" json:"lock_whitelist_id"`
	LockUseragentId         uint     `form:"lock_useragent_id" json:"lock_useragent_id"`
	LockWhitelistReverse    uint     `form:"lock_whitelist_reverse" json:"lock_whitelist_reverse"`
	LockUseragentReverse    uint     `form:"lock_useragent_reverse" json:"lock_useragent_reverse"`
	Paused                  bool     `form:"paused" json:"paused"`
	AutoRefresh             bool     `form:"auto_refresh" json:"auto_refresh"`
	StripReferrer           bool     `form:"strip_referrer" json:"strip_referrer"`
	ForceRefresh            bool     `form:"force_refresh" json:"force_refresh"`
	StripQueryString        bool     `form:"strip_query_string" json:"strip_query_string"`
	ForceHost               uint     `form:"force_host" json:"force_host"`
	AllowedQueryStrings     []string `form:"allowed_qs[]" sql:"-" json:"allowed_qs"`
	BbsiPath                string   `form:"bbsi_path" json:"bbsi_path"`
	SortMethod              uint     `form:"sort_method" json:"sort_method"`
	ScoringTimeout          uint     `form:"scoring_timeout" json:"scoring_timeout"`
	ScoringRedirectEnabled  bool     `form:"scoring_redirect_enabled" json:"scoring_redirect_enabled"`
	ScoringRedirectOverride string   `form:"scoring_redirect_override" json:"scoring_redirect_override"`
	Campaigns               []AdscoopRedirectCampaignRead
	BustIframe              bool
}

type AdscoopTracking struct {
	Timeslice        time.Time `sql:"unique_index:timeslice"`
	RedirectID       uint      `sql:"unique_index:timeslice"`
	UrlId            uint      `sql:"unique_index:timeslice"`
	Cpc              string    `sql:"type:decimal(11,6);unique_index:timeslice"`
	UniqueIdentifier string    `sql:"unique_index:timeslice"`
	Count            uint
	Engagement       uint
	Load             uint
	TimeOnSite       string
	TimeOnSiteCount  string
}

type AdscoopCampaign struct {
	gorm.Model
	ClientID                    uint                  `form:"client_id" sql:"index"`
	Cpc                         string                `sql:"type:decimal(11,6)" form:"cpc"`
	Name                        string                `form:"name"`
	Paused                      bool                  `form:"paused"`
	DailyImpsLimit              uint                  `form:"daily_imps_limit"`
	TrackingMethod              uint                  `form:"tracking_method"`
	Source                      string                `form:"source"`
	Type                        uint                  `form:"type"`
	StartDatetimeEdit           string                `form:"start_datetime_edit" sql:"-"`
	EndDatetimeEdit             string                `form:"end_datetime_edit" sql:"-"`
	XmlType                     uint                  `form:"xml_type"`
	XmlUrl                      string                `form:"xml_url" sql:"-"`
	PerformanceBasedPauseEnable bool                  `form:"performance_based_pause_enable"`
	PerformanceBasedPauseResume bool                  `form:"performance_based_pause_resume"`
	PerformanceBasedCompareA    uint                  `form:"performance_based_compare_a"`
	PerformanceBasedCompareB    uint                  `form:"performance_based_compare_b"`
	PerformanceBasedPercent     uint                  `form:"performance_based_percent"`
	PerformanceBasedPauseQueued uint                  `form:"performance_based_pause_queued"`
	PerformanceBasedNotifyOnly  bool                  `form:"performance_based_notify_only"`
	EnableCampaignQualityCheck  bool                  `form:"enable_campaign_quality_check"`
	StartDatetime               time.Time             `form:"-"`
	EndDatetime                 time.Time             `form:"-"`
	DisableStartTime            bool                  `form:"disable_start_time"`
	DisableEndTime              bool                  `form:"disable_end_time"`
	WeightVariance              uint                  `form:"weight_variance"`
	EnableStartStopTimes        bool                  `form:"enable_start_stop_times"`
	WeightsLastUpdated          time.Time             `form:"-"`
	Urls                        []*AdscoopCampaignUrl `form:"-" sql:"-" json:"urls"`
	AllUrls                     []*AdscoopCampaignUrl `form:"-" sql:"-" json:"all_urls"`
	ActiveUrls                  []string              `form:"url[]" sql:"-"`
	ActiveWeights               []uint                `form:"weight[]" sql:"-"`
	Inactive                    bool                  `form:"inactive"`
	EnableUnloadTracking        bool                  `form:"enable_unload_tracking"`
	IsRon                       bool
	CampaignGroupId             uint `form:"campaign_group_id"`
	CampaignGroupWeight         uint `form:"campaign_group_weight"`
	AppendRc                    bool `form:"append_rc"`
}

type AdscoopCampaignGroup struct {
	Id   uint   `form:"id"`
	Name string `form:"name"`
}

type AdscoopCampaignScheduleAddons struct {
	ScheduleExecuteEdit string    `sql:"-" form:"schedule_execute_edit"`
	ScheduleExecute     time.Time `form:"-"`
	CampaignID          uint      `form:"campaign_id" sql:"index"`
	ScheduleLabel       string    `form:"schedule_label"`
	ScheduleQueued      bool      `form:"schedule_queued"`
	SchedulePending     bool      `form:"schedule_pending"`
}

type AdscoopCampaignSchedule struct {
	AdscoopCampaign
	AdscoopCampaignScheduleAddons
}

func (a AdscoopCampaignSchedule) TableName() string {
	return "adscoop_campaign_schedules"
}

type AdscoopCampaignScheduleUrl struct {
	AdscoopCampaignUrl
}

func (a AdscoopCampaignScheduleUrl) TableName() string {
	return "adscoop_scheduled_urls"
}

type AdscoopCampaignUrl struct {
	gorm.Model
	AdscoopCampaignID uint `gorm:"column:campaign_id" sql:"unique_index:campaign_url"`
	Weight            uint
	Url               string `sql:"unique_index:campaign_url"`
	Title             string
}

func (a AdscoopCampaignUrl) TableName() string {
	return "adscoop_urls"
}

type AdscoopWhitelist struct {
	ID   uint     `form:"id" gorm:"unique_index"`
	Name string   `form:"name"`
	Urls []string `form:"url[]" sql:"-"`
}

type AdscoopWhitelistUrl struct {
	ID                 uint   `form:"id" gorm:"primary_key"`
	Url                string `form:"url" sql:"unique_index:whitelist_url"`
	AdscoopWhitelistId uint   `form:"whitelist_id" sql:"unique_index:whitelist_url"`
}

type AdscoopWhitelistUseragentGroup struct {
	ID         uint     `form:"id" gorm:"primary_key"`
	Name       string   `form:"name"`
	Useragents []string `form:"ua[]" sql:"-"`
}

type AdscoopWhitelistUseragent struct {
	ID                               uint   `form:"id"`
	Useragent                        string `form:"ua" sql:"unique_index:whitelist_ua"`
	AdscoopWhitelistUseragentGroupId uint   `form:"whitelist_useragent_group" sql:"unique_index:whitelist_ua"`
}

type RssXml struct {
	XMLName xml.Name      `xml:"rss"`
	Ver     string        `xml:"version,attr"`
	Channel RssXmlChannel `xml:"channel"`
	Media   string        `xml:"xmlns:media,attr"`
	Key     string        `xml:"-"`
	Config  string        `xml:"-"`
	Host    string        `xml:"-"`
	Cached  bool          `xml:"-"`
}

type RssXmlChannel struct {
	XMLName     xml.Name     `xml:"channel"`
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	Link        string       `xml:"link"`
	Items       []RssXmlItem `xml:"item"`
}

type RssXmlItem struct {
	XMLName xml.Name `xml:"item"`
	Link    string   `xml:"link"`
	Guid    struct {
		Guid        string `xml:",chardata"`
		IsPermalink bool   `xml:"isPermalink,attr"`
	} `xml:"guid"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Media       string `xml:"http://search.yahoo.com/mrss media"`
	Thumbnail   struct {
		XMLName xml.Name `xml:"media:thumbnail"`
		Url     string   `xml:"url,attr"`
	}
	PubDate string `xml:"pubDate"`
}

type SQLDatabase struct {
	DB gorm.DB
}

type RedirectManager struct {
	SQLDatabase
}

type CampaignManager struct {
	SQLDatabase
}

func (db *SQLDatabase) Connect(env, dbname string) (err error) {
	/*dbConn := fmt.Sprintf(os.Getenv(env), dbname)
	dbConn += "&charset=utf8"

	db.DB, err = gorm.Open("mysql", dbConn)*/
	return
}

func (db *SQLDatabase) GenerateTables() (err error) {
	err = db.DB.AutoMigrate(&AdscoopCampaign{},
		&AdscoopClientEmail{}, &AdscoopClientTransaction{}, &AdscoopClient{},
		&AdscoopFeedRedirect{}, &AdscoopFeed{}, &AdscoopHost{}, &AdscoopPaymentHash{},
		&AdscoopRedirectCampaign{}, &AdscoopRedirectQuerystring{},
		&AdscoopRedirect{}, &AdscoopTracking{}, &AdscoopCampaignUrl{},
		&AdscoopWhitelistUrl{}, &AdscoopWhitelistUseragent{}, &AdscoopWhitelistUseragentGroup{},
		&AdscoopWhitelist{}, &AdscoopPastebin{}).Error
	return
}
func (db *SQLDatabase) DeleteTables() (err error) {
	err = db.DB.DropTable(&AdscoopCampaign{},
		&AdscoopClientEmail{}, &AdscoopClientTransaction{}, &AdscoopClient{},
		&AdscoopFeedRedirect{}, &AdscoopFeed{}, &AdscoopHost{}, &AdscoopPaymentHash{},
		&AdscoopRedirectCampaign{}, &AdscoopRedirectQuerystring{},
		&AdscoopRedirect{}, &AdscoopTracking{}, &AdscoopCampaignUrl{},
		&AdscoopWhitelistUrl{}, &AdscoopWhitelistUseragent{}, &AdscoopWhitelistUseragentGroup{},
		&AdscoopWhitelist{}, &AdscoopPastebin{}).Error
	return
}

type AdscoopStats struct {
	ClientID          uint
	ClientName        string
	TotalSpend        string
	TodaySpend        string
	CampaignCount     int64
	ChargeAmount      string
	ChargeAmountFloat float64
	Balance           string
	TotalImpressions  string
	Campaigns         []AdscoopStatsCampaign
}

type AdscoopStatsCampaign struct {
	CampaignID             int64
	CampaignName           string
	CampaignSource         string
	CampaignGoal           int
	CampaignGoalString     string
	CampaignCpc            string
	TotalSpend             float64
	TotalSpendString       string
	TotalImpressions       int
	TotalImpressionsString string
	TotalVerified          int
	TotalVerifiedString    string
	TotalLoad              int
	TotalLoadString        string
	TodayImpressions       int
	TodayVerified          int
	TodayVerifiedString    string
	TodayLoad              int
	TodayLoadString        string
	TodayImpressionsString string
	TodaySpend             float64
	TodaySpendString       string
	DailyImpsLimit         int
	IsDone                 bool
	Paused                 int64
	TrackingMethod         int64
	TimeOnSite             string
	Redirs                 []AdscoopRedirect
}

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

type AdscoopPastebin struct {
	gorm.Model
	Title string        `form:"title"`
	Text  string        `sql:"type:longtext" form:"text"`
	HTML  template.HTML `sql:"-"`
}
