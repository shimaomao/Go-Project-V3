package adscoopsCaches

import (
	"encoding/json"
	"fmt"

	"app/structs"
)

func LoadAdscoopQueryStrings(id uint) (asqs []structs.RedirectQuerystring) {
	// log.Println("looking for query strings in DB")
	AdscoopsDB.Where("redirect_id = ?", id).
		Find(&asqs)

	b, err := json.Marshal(asqs)

	if err == nil {
		rp := RedisPool.Get()
		defer rp.Close()
		rp.Do("SET", fmt.Sprintf(QUERYSTRING_KEY, id), b)
	}

	return
}
