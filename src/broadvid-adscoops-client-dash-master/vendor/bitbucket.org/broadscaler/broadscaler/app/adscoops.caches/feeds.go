package adscoopsCaches

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/broadscaler/broadscaler/app/structs"
)

func LoadAdscoopFeeds(hash string) (asf structs.AdscoopFeed) {

	AdscoopsDB.Where("paused = 0 AND hash = ?", hash).
		Find(&asf)

	AdscoopsDB.Where("adscoop_feed_redirects.feed_id = ?", asf.Id).Joins(`JOIN adscoop_feed_redirects ON adscoop_feed_redirects.redirect_id = adscoop_redirects.id`).
		Find(&asf.Redirects)

	b, err := json.Marshal(asf)

	if err == nil {
		rp := RedisPool.Get()
		rp.Do("SET", fmt.Sprintf(FEEDS_KEY, hash), b)
	}
	return
}
