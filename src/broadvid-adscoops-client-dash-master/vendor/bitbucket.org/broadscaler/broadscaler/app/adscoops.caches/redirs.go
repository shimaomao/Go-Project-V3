package adscoopsCaches

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/structs"
)

func updateRedirCaches() {
	var campaigns structs.Campaigns
	campaigns.UpdateTodayCount()
	var redirects structs.Redirects
	err := redirects.GetRecent()

	if err != nil {
		log.Errorf("Error getting recent redirs: %s", err)
		return
	}

	for _, r := range redirects {
		UpdateRedirect(&r)
	}
}

func UpdateRedirect(r *structs.Redirect) {
	rc := RedisPool.Get()
	defer rc.Close()

	var psdata structs.PubSubMessage
	psdata.Event = AdscoopPSEventRedirect
	psdata.Data.IDInt = r.ID
	psdata.Data.IDString = r.Hash

	if err := r.CheckDailyProgress(); err != nil {
		log.Errorf("Couldn't check daily progress: %s", err)
	}

	LoadAdscoopRedirs(r.Hash)
	LoadAdscoopQueryStrings(r.ID)
	LoadAdscoopRedirsCampaign(r.ID)

	rc.Do("PUBLISH", AdscoopRedisPubSubChannel, psdata.JSONify())
}

func LoadAdscoopRedirs(hash string) (asc structs.Redirect) {

	AdscoopsDB.Where("paused = 0 AND hash = ?", hash).
		Find(&asc)

	if asc.ForceHost != "0" {
		var h structs.Host
		AdscoopsDB.Find(&h, asc.ForceHost)
		if h.ID != 0 {
			asc.ForceHostString = h.Host
		}
	}

	if asc.LockWhitelistId != "0" {
		AdscoopsDB.Where("adscoop_whitelist_id = ?", asc.LockWhitelistId).Find(&asc.LockWhitelistUrls)
	}

	if asc.LockUseragentId != "0" {
		AdscoopsDB.Where("adscoop_whitelist_useragent_group_id = ?", asc.LockUseragentId).Find(&asc.LockUseragents)
	}

	b, err := json.Marshal(asc)

	if err == nil {
		rp := RedisPool.Get()
		defer rp.Close()
		rp.Do("SET", fmt.Sprintf(REDIRS_KEY, hash), b)
	}

	return
}
