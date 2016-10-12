package app

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var TimeLayout = "2006-01-02 15:04:05"

func adscoopsClientsCampaignsScheduleNew(r render.Render, p martini.Params, user *UserWithPolicy, req *http.Request) {
	var retData struct {
		DefaultRetData
		AdscoopClient         structs.AdscoopClient
		AdscoopCampaignGroups []structs.AdscoopCampaignGroup
		AdscoopCampaign       struct {
			structs.AdscoopCampaign
			structs.AdscoopCampaignScheduleAddons
		}
		UserPolicy structs.AdscoopUserPolicy
	}
	retData.User = user

	var client structs.AdscoopClient
	db.Find(&client, user.ClientID)

	var up structs.AdscoopUserPolicy
	db.Find(&up, user.UserPolicyID)

	if p["cid"] == "" && up.AddCopyCampaignEnabled == false {
		r.Redirect("/")
		return
	}

	location, _ := time.LoadLocation("America/Los_Angeles")

	db.Find(&retData.AdscoopClient, user.ClientID)

	if retData.AdscoopClient.ID == 0 {
		return
	}

	db.Find(&retData.UserPolicy, user.UserPolicyID)

	db.LogMode(true)
	db.Select("adscoop_campaign_groups.*").Table("adscoop_campaign_groups").Joins("JOIN adscoop_client_campaign_groups ON adscoop_client_campaign_groups.campaign_group_id = adscoop_campaign_groups.id").Where("adscoop_client_campaign_groups.client_id = ?", client.ID).Find(&retData.AdscoopCampaignGroups)
	db.LogMode(false)

	if p["sid"] == "" {

		if p["cid"] == "" {
			// Creating a new campaign

			log.Println("THIS IS A NEW CAMPAIGN")

			retData.AdscoopCampaign.ID = 0
			retData.AdscoopCampaign.Type = 0
			retData.AdscoopCampaign.Cpc = client.DefaultCpc
		} else {
			db.Table("adscoop_campaigns").Find(&retData.AdscoopCampaign, p["cid"])
			if req.FormValue("copy") == "true" && up.AddCopyCampaignEnabled {
				// Copying existing new campaign

				retData.AdscoopCampaign.ID = 0
			}
		}

		retData.AdscoopCampaign.CampaignID = retData.AdscoopCampaign.ID
		retData.AdscoopCampaign.ID = 0
		log.Println("now", time.Now())
		retData.AdscoopCampaign.ScheduleExecute = time.Now().UTC()
		// retData.AdscoopCampaign.ScheduleExecute.In(location)
	} else {
		db.Table("adscoop_campaign_schedules").Find(&retData.AdscoopCampaign, p["sid"])
	}

	loc := time.Date(retData.AdscoopCampaign.ScheduleExecute.Year(),
		retData.AdscoopCampaign.ScheduleExecute.Month(),
		retData.AdscoopCampaign.ScheduleExecute.Day(),
		retData.AdscoopCampaign.ScheduleExecute.Hour(),
		retData.AdscoopCampaign.ScheduleExecute.Minute(), 0, 0, time.UTC)

	retData.AdscoopCampaign.ScheduleExecute = loc.In(location)

	if retData.AdscoopCampaign.Type == 0 {
		var us []*structs.AdscoopCampaignUrl

		if p["sid"] == "" {
			db.Where("campaign_id = ?", p["cid"]).Order("id desc").Find(&us)
		} else {
			db.Table("adscoop_scheduled_urls").Where("campaign_id = ?", p["sid"]).Order("id desc").Find(&us)
		}

		for _, y := range us {
			retData.AdscoopCampaign.Urls = append(retData.AdscoopCampaign.Urls, y)
		}

		us = []*structs.AdscoopCampaignUrl{}
		if p["sid"] == "" {
			db.Unscoped().Where("deleted_at != '0000-00-00' AND campaign_id = ?", p["cid"]).Order("id desc").Find(&us)
		} else {
			db.Table("adscoop_scheduled_urls").Unscoped().Where("deleted_at != '0000-00-00' AND campaign_id = ?", p["sid"]).Order("id desc").Find(&us)
		}

		for _, y := range us {
			retData.AdscoopCampaign.AllUrls = append(retData.AdscoopCampaign.AllUrls, y)
		}

		r.HTML(http.StatusOK, "adscoops/clients/campaigns/new_schedule", retData)
	}

	if retData.AdscoopCampaign.Type == 1 || retData.AdscoopCampaign.Type == 2 {
		var asu structs.AdscoopCampaignUrl

		if p["sid"] == "" {
			db.Where("campaign_id = ? AND weight = 0", p["cid"]).Find(&asu)
		} else {
			db.Table("adscoop_scheduled_urls").Where("campaign_id = ? AND weight = 0", p["sid"]).Find(&asu)
		}
		retData.AdscoopCampaign.XmlUrl = asu.Url
		if retData.AdscoopCampaign.Type == 1 {
			r.HTML(http.StatusOK, "adscoops/clients/campaigns/new_xml_schedule", retData)
		}

		if retData.AdscoopCampaign.Type == 2 {
			r.HTML(http.StatusOK, "adscoops/clients/campaigns/new_xml_ingest_schedule", retData)
		}
	}
}

func adscoopsClientsCampaignsScheduleSave(user *UserWithPolicy, r render.Render, asc structs.AdscoopCampaignSchedule, req *http.Request) {
	iid, err := strconv.Atoi(req.FormValue("ID"))
	if err != nil {
		r.Text(500, "Cannot save campaign at this time")
		return
	}
	asc.ID = uint(iid)
	if asc.Type == 0 {
		if len(asc.ActiveUrls) == 0 {
			asc.Paused = true
		}
	}

	var client structs.AdscoopClient

	db.Find(&client, user.ClientID)

	if asc.CampaignID == 0 {
		var newCampaign structs.AdscoopCampaign
		newCampaign.Name = asc.Name
		newCampaign.Type = newCampaign.Type
		newCampaign.ClientID = user.ClientID
		newCampaign.Cpc = client.DefaultCpc
		db.Save(&newCampaign)

		asc.CampaignID = newCampaign.ID
	}

	layout := "01/02/2006 3:04 PM"

	asc.StartDatetime, err = time.Parse(layout, asc.StartDatetimeEdit)
	if err != nil {
		log.Println("start err", err)
	}
	asc.EndDatetime, err = time.Parse(layout, asc.EndDatetimeEdit)
	if err != nil {
		log.Println("end err", err)
	}

	location, _ := time.LoadLocation("America/Los_Angeles")

	today := time.Now()

	today = today.In(location)

	asc.ScheduleExecute, err = time.Parse(layout, asc.ScheduleExecuteEdit)
	if err != nil {
		log.Println("execute err", err)
	}

	today = time.Date(asc.ScheduleExecute.Year(), asc.ScheduleExecute.Month(), asc.ScheduleExecute.Day(), asc.ScheduleExecute.Hour(),
		asc.ScheduleExecute.Minute(), 0, 0, location)

	today = today.In(time.UTC)

	log.Println("today", today)

	asc.ScheduleExecute = today

	asc.PerformanceBasedPauseQueued = 1
	asc.SchedulePending = true
	asc.ScheduleQueued = false
	asc.ClientID = user.ClientID
	err = db.Where("client_id = ?", user.ClientID).Save(&asc).Error

	rand.Seed(time.Now().Unix())

	if req.FormValue("macro_replace") != "" {
		log.Println("replace", req.FormValue("macro_replace"))
		find := req.FormValue("macro_find")

		if find == "" {
			find = "[REPLACE_ME]"
		}
		for x, y := range asc.ActiveUrls {
			asc.ActiveUrls[x] = strings.Replace(y, find, req.FormValue("macro_replace"), -1)
		}
	}

	if asc.Type == 0 {
		db.Where("campaign_id = ?", asc.ID).Delete(structs.AdscoopCampaignScheduleUrl{})
		for x, y := range asc.ActiveUrls {
			var asu structs.AdscoopCampaignScheduleUrl
			asu.Url = y
			var urlWeight uint
			if asc.WeightVariance != 0 {
				urlWeight = uint(rand.Intn(100-(100-int(asc.WeightVariance))) + (100 - int(asc.WeightVariance)))
			} else {
				urlWeight = asc.ActiveWeights[x]
			}
			asu.Weight = urlWeight
			asu.AdscoopCampaignID = asc.ID
			info := db.Save(&asu)

			if info.Error != nil {
				db.Where("campaign_id = ? AND url = ?", asc.ID, y).Find(&asu)
				if asu.ID == 0 {
					_, err := db.DB().Exec("UPDATE adscoop_scheduled_urls SET deleted_at = '0000-00-00', weight = ? WHERE campaign_id = ? AND url = ?", urlWeight, asc.ID, y)
					log.Println("url update err", err)
				}
			}
		}
	}

	if asc.Type == 1 || asc.Type == 2 {
		var asu structs.AdscoopCampaignScheduleUrl
		log.Println("asc", asc.ID)
		db.Where("campaign_id = ? AND weight = 0", asc.ID).Find(&asu)
		log.Println("url", asc.XmlUrl)
		log.Println("asu", asu.ID)
		asu.Url = asc.XmlUrl
		asu.AdscoopCampaignID = asc.ID
		asu.Weight = 0
		db.Save(&asu)
		if asc.Type == 2 {
			asUtils.IngestXml(asc.ID)
		}
	}

	go createCampaignCache(user.ClientID)

	r.Redirect("/")
}
