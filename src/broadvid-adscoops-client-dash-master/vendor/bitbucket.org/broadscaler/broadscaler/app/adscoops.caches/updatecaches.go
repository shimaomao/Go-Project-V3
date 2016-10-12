package adscoopsCaches

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bradhe/stopwatch"
)

func UpdateCaches() {
	start := stopwatch.Start()
	log.Println("updating Caches")

	updateCampaignUrlCaches()
	updateRedirCaches()
	// updateFeedCaches() // TODO: build out feed support

	watch := stopwatch.Stop(start)

	log.Printf("done updating caches, took: %vms", watch.Milliseconds())
}

func BeginCacheUpdater() {
	InitStructs()
	for {
		UpdateCaches()
		time.Sleep(10 * time.Second)
	}
}
