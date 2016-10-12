package frontendControllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/patrickmn/go-cache"
	"github.com/unrolled/render"
	"github.com/jinzhu/gorm"
	"app/adscoops.caches"
	"app/structs"

)

var (
	db                  *gorm.DB
	gc               *cache.Cache
	redisPool        *redis.Pool
)


func setup() (success bool) {
	var err error
	//adscoopsDB := "root:@/adscoops?parseTime=true"
	//db, err = gorm.Open("mysql", adscoopsDB)
	db =structs.AdscoopsDB
	if(err!=nil){
		panic(err.Error())
	}
	adscoopsCaches.InitStructs()
	go tsc.Save()

	go startCachesTimer()

	success = true
	return
}


/*Setup is to set up the controllers for the frontend for adscoops and broadvids*/
func Setup(m *gin.Engine) {
	log.Println("Setting up frontend controllers...")


	urlTrackingMethod.list = make(map[uint]string)
	// Setting up to listen to when the application ends, we are going to force push stats.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		<-ch
		td.Push()
		os.Exit(1)
	}()

	// Set up tracking stats list
	td.list = make(map[string]trackingData)

	// Set up tracking stats ticker to save data every minute
	go td.Save()

	gc = cache.New(60*time.Minute, 30*time.Second)

	setup()


	r = render.New(render.Options{
		IsDevelopment: false,
		PrefixXML:     []byte("<?xml version='1.0' encoding='UTF-8'?>"), // Prefixes XML responses with the given bytes.
	})



	m.LoadHTMLGlob("./public/views/templates/**/*")
	log.Println("we are ruccccccnning frontend")
	m.GET("/tracking.js",controllerTrackingJs)
	m.GET("/loadClient", controllerLoadClient)
	m.GET("/f/{hash}.{type}", controllerFeed)
	m.GET("/r/{hash}", controllerRedirect)
	m.GET("/lastUpdated/{hash}", controllerLastUpdated)
	m.GET("/v/{hash}", controllerValidUser)
	m.GET("/u/{hash}", controllerRedirectUrl)
	m.GET("/t/engagement", controllerTrackEngagement)
	m.GET("/t/load", controllerTrackLoad)
	m.GET("/t/tos", controllerTrackTimeOnSite)

}