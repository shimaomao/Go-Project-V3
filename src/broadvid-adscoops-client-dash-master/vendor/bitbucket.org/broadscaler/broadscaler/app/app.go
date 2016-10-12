package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/patrickmn/go-cache"

	"github.com/jinzhu/gorm"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool
var Cache *cache.Cache

const (
	RedisAuthPassword            = "broadvid123!!!"
	AdscoopTempStatsRedisChannel = "adscoops-temp-stats"
	// googleApiKey                 = "AIzaSyCrXH-owMYco_PB3mI30Tkp2m9yhtsP18M"
	// googleCallback               = "http://localhost:3000/oauth2callback"
	googleApiKey   = "AIzaSyAyuE9xg-pgGxejq8R9U9_5dXt1_Ih-cx8"
	googleCallback = "http://localhost:8080/oauth2callback"
)

func init() {
	log.Println("Broadscaler app initializing...")
}
