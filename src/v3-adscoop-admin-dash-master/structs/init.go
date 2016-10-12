package structs

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

var (
	AdscoopsDB         *gorm.DB
	AdscoopsRealtimeDB *gorm.DB
	BroadvidDB         *gorm.DB
	RedisPool          *redis.Pool
	Cache              *cache.Cache
)

func init() {
	Cache = cache.New(time.Hour, 30*time.Second)
}

const (
	layout                          = "01/02/2006 3:04 PM"
	lookbackInMinutes time.Duration = 100
)

type GroupIntf interface {
	FindAll() error
}

type GroupVisbleIntf interface {
	FindVisible(userID uint) error
}

type GroupScheduleIntf interface {
	FindAll(id string) error
}

type SingleIntf interface {
	Find(id string) error
	Save() error
}

type SingleIntfByRedir interface {
	Find(id string) error
	Save(id string) error
}

type BasicIntf interface {
	BasicSave() error
}
