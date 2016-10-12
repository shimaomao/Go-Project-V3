package adminControllers

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
)

var AdscoopsDB *gorm.DB
var AdscoopsRealtimeDB *gorm.DB
var BroadvidDB *gorm.DB
var RedisPool *redis.Pool

func init() {

	host := "http://localhost:8080"

	if os.Getenv("PRODUCTION_HOST") != "" {
		host = os.Getenv("PRODUCTION_HOST")
	}

	redirectUrl := fmt.Sprintf("%s/auth/callback?provider=gplus", host)
	gothic.Store = sessions.NewFilesystemStore(os.TempDir(), []byte("auth-sess"))

	goth.UseProviders(
		// gplus.New("462194986011-vj5orge2b7o0rthrji7rff89m5866u47.apps.googleusercontent.com", "xRTZL02LEZaHIukmEae6xq60", redirectUrl),
		gplus.New("772951475452-31kbfmrn1ku8mnjh3v8520teekpg0lr6.apps.googleusercontent.com", "TRSIe9ZmzMieJNQ-nPd4Kmrb", redirectUrl),
	)
}
