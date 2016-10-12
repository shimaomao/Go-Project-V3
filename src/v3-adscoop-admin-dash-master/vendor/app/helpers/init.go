package helpers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool
