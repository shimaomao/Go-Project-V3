package adscoopsCaches

import (
	"app/adonnetwork"
	"app/structs"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

const (
	QUERYSTRING_KEY            = "adscoopRedirsQueryString_%v"
	REDIRS_KEY                 = "adscoopRedirs_%s"
	FEEDS_KEY                  = "adscoopFeeds_%s"
	CAMPAIGNS_KEY              = "adscoopRedirsCampaign_%v"
	URLS_KEY                   = "adscoopCampaignUrls_%v"
	HOSTS_BY_ID_KEY            = "adscoopHostById_%s"
	HOSTS_BY_HOST_KEY          = "adscoopHostByHost_%s"
	AdscoopPSEventRedirect     = "redirect"
	AdscoopPSEventCampaignUrls = "campaign-urls"
	AdscoopRedisPubSubChannel  = "adscoops-updates"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool
var Cache *cache.Cache

func InitStructs() {
	structs.AdscoopsDB = AdscoopsDB
	structs.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	structs.BroadvidDB = BroadvidDB
	structs.RedisPool = RedisPool
	structs.Cache = Cache
	adonnetwork.AdscoopsDB = AdscoopsDB
}
