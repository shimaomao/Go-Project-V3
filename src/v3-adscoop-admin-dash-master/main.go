package main

import (
	"app"
	"app/configSettting"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// broadscaler admin
	//
	// adscoopsDB := fmt.Sprintf(os.Getenv("GO_DATABASE_CONN_ADSCOOPS"), "adscoops")
	// adscoopsRealtimeDB := fmt.Sprintf(os.Getenv("GO_DATABASE_CONN_ADSCOOPS_RT"), "adscoops")
	//
	// broadvidadsDB := fmt.Sprintf(os.Getenv("GO_DATABASE_CONN"), "broadvidadserver")

	adscoopsDB := "root:@/adscoops?parseTime=true"
	adscoopsRealtimeDB := "root:@/adscoops?parseTime=true"
	broadvidadsDB := "root:@/broadvidadserver?parseTime=true"

	//Setup Adscoops DB connection
	db, err := gorm.Open("mysql", adscoopsDB)

	if err != nil {
		log.Panicf("Error connecting to the DB: %s", err)
	}
	db.LogMode(true)

	log.Info("Connected to the first DB connection")

	//Setup Adscoops RealTime DB connection

	dbrt, err := gorm.Open("mysql", adscoopsRealtimeDB)

	if err != nil {
		log.Panicf("Error connecting to the DB: %s", err)
	}
	dbrt.LogMode(true)

	log.Info("Connected to the second DB connection")

	//Setup Broadvid DB connection

	dbbv, err := gorm.Open("mysql", broadvidadsDB)

	if err != nil {
		log.Panicf("Error connecting to the DB: %s", err)
	}

	log.Info("Connected to the first broadvid DB connection")

	// Setup Redis

	redisHost := "localhost:6379"

	if os.Getenv("REDIS_HOST") != "" {
		redisHost = os.Getenv("REDIS_HOST")
	}

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		var c redis.Conn
		var err error
		c, err = redis.Dial("tcp", redisHost)

		if err != nil {
			log.Panic("REDIS err", err)
			return nil, err
		}

		if os.Getenv("REDIS_HOST") != "" {
			c.Do("AUTH", configSettting.RedisAuthPassword)

		}

		return c, err
	}, 10)

	log.Info("Connected to redis")
	log.Println(redisPool)
	log.Println(dbbv)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	configSettting.AdscoopsDB = db
	configSettting.AdscoopsRealtimeDB = dbrt
	configSettting.BroadvidDB = dbbv
	configSettting.RedisPool = redisPool
	app.InitAdmin()
}
