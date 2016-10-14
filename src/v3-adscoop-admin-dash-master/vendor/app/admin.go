package app

import (
	"encoding/json"
	"os"
	"strings"

	// "os"
	// "strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"app/api.controllers"
	"github.com/garyburd/redigo/redis"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"app/adscoops.caches"
	"app/configSettting"
	"app/frontend.controllers"
	"app/janitor"
	"app/payment.controllers"
	"app/structs"
)

var tsc structs.TempStatsContainer

func redisListenForAdscoopStats() {
	go tsc.Save()
	c := configSettting.RedisPool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe(configSettting.AdscoopTempStatsRedisChannel)

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
	configSettting.AdscoopsDB.LogMode(true)
	configSettting.AdscoopsRealtimeDB.LogMode(true)

	go redisListenForAdscoopStats()

	m := gin.Default()

	store := ginsessions.NewCookieStore([]byte("hello123!!!"))
	m.Use(ginsessions.Sessions("broadscaler", store))

	apiControllers.AdscoopsDB = configSettting.AdscoopsDB
	apiControllers.AdscoopsRealtimeDB = configSettting.AdscoopsRealtimeDB
	apiControllers.BroadvidDB = configSettting.BroadvidDB
	apiControllers.RedisPool = configSettting.RedisPool

	apiControllers.Setup(m)
	paymentController.Setup(m)
	frontendControllers.Setup(m)

	adscoopsCaches.AdscoopsDB = configSettting.AdscoopsDB
	adscoopsCaches.AdscoopsRealtimeDB = configSettting.AdscoopsRealtimeDB
	adscoopsCaches.BroadvidDB = configSettting.BroadvidDB
	adscoopsCaches.RedisPool = configSettting.RedisPool
	adscoopsCaches.Cache = cache.New(1*time.Minute, 30*time.Second)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) != "development" {
		go adscoopsCaches.BeginCacheUpdater()

		janitor.AdscoopsDB = configSettting.AdscoopsDB
		go janitor.Run()
	}
	m.Run()

}
