package app

import (
	"encoding/json"
	"os"
	"strings"
	// "os"
	// "strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/admin.controllers"
	"bitbucket.org/broadscaler/broadscaler/app/adscoops.caches"
	"bitbucket.org/broadscaler/broadscaler/app/janitor"
	// "bitbucket.org/broadscaler/broadscaler/app/janitor"
	"bitbucket.org/broadscaler/broadscaler/app/structs"

	"github.com/garyburd/redigo/redis"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var tsc structs.TempStatsContainer

func redisListenForAdscoopStats() {
	go tsc.Save()
	c := RedisPool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe(AdscoopTempStatsRedisChannel)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var data structs.TempStats
			err := json.Unmarshal(v.Data, &data)

			log.Printf("received: %s", v.Data)

			if err == nil {
				tsc.Add(&data)
			}
		}
	}
}

/*InitAdmin is the bootstrapping of the admin application*/
func InitAdmin() {
	AdscoopsDB.LogMode(true)
	AdscoopsRealtimeDB.LogMode(true)

	go redisListenForAdscoopStats()

	m := gin.Default()

	store := ginsessions.NewCookieStore([]byte("hello123!!!"))
	m.Use(ginsessions.Sessions("broadscaler", store))

	adminControllers.AdscoopsDB = AdscoopsDB
	adminControllers.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	adminControllers.BroadvidDB = BroadvidDB
	adminControllers.RedisPool = RedisPool

	adminControllers.Setup(m)

	adscoopsCaches.AdscoopsDB = AdscoopsDB
	adscoopsCaches.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	adscoopsCaches.BroadvidDB = BroadvidDB
	adscoopsCaches.RedisPool = RedisPool
	adscoopsCaches.Cache = cache.New(1*time.Minute, 30*time.Second)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) != "development" {
		go adscoopsCaches.BeginCacheUpdater()

		janitor.AdscoopsDB = AdscoopsDB
		go janitor.Run()
	}
	m.Run()

}
