package configSettting

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/patrickmn/go-cache"
	//"structs"
	//"github.com/martini-contrib/render"
	//"github.com/unrolled/render"
	"github.com/jinzhu/gorm"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool
var Cache *cache.Cache

/*var TSC structs.TempStatsContainer
var Ren *render.Render
var TD  =structs.ClickStats{}*/

const (
	RedisAuthPassword            = "broadvid123!!!"
	AdscoopTempStatsRedisChannel = "adscoops-temp-stats"
	// googleApiKey                 = "AIzaSyCrXH-owMYco_PB3mI30Tkp2m9yhtsP18M"
	// googleCallback               = "http://localhost:3000/oauth2callback"
	googleApiKey   = "AIzaSyAyuE9xg-pgGxejq8R9U9_5dXt1_Ih-cx8"
	googleCallback = "http://localhost:8080/oauth2callback"

	// frontend contstants

	redirContext      = "adscoopRedirRetData"
	QUERYSTRING_KEY   = "adscoopRedirsQueryString_%v"
	REDIRS_KEY        = "adscoopRedirs_%s"
	FEEDS_KEY         = "adscoopFeeds_%s"
	CAMPAIGNS_KEY     = "adscoopRedirsCampaign_%v"
	URLS_KEY          = "adscoopCampaignUrls_%v"
	HOSTS_BY_ID_KEY   = "adscoopHostById_%s"
	HOSTS_BY_HOST_KEY = "adscoopHostByHost_%s"
)

func init() {
	log.Println("Broadscaler app initializing...")
}
