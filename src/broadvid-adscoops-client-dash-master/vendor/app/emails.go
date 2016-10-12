package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
	"github.com/mailgun/mailgun-go"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sendEmailForClient(user *UserWithPolicy) {

	var retData struct {
		State string
		DefaultRetData
		Client             structs.AdscoopClient
		ScheduledCampaigns []struct {
			structs.AdscoopCampaignSchedule
			Comparisons []template.HTML
		}
	}

	var up structs.AdscoopUserPolicy
	db.Find(&up, user.UserPolicyID)
	retData.State = "Approved"
	retData.User = user
	db.Find(&retData.Client, user.ClientID)

	db.Select("*").Table("adscoop_campaign_schedules").Where("adscoop_campaign_schedules.schedule_pending > 0 AND adscoop_campaign_schedules.client_id = ?", user.ClientID).Find(&retData.ScheduledCampaigns)

	for i := range retData.ScheduledCampaigns {
		c := retData.ScheduledCampaigns[i]

		var exc structs.AdscoopCampaign
		db.Find(&exc, c.CampaignID)

		c.Comparisons = append(c.Comparisons, template.HTML("<tr style='border-bottom: 1px solid black; background-color:#000;'><td style='padding:1em; color:#fff !important; font-weight:800;'>Field</td><td style='padding:1em; color:#fff !important; font-weight:800;'>Old</td><td style='padding:1em; color:#fff !important; font-weight:800;'>New</td></tr>"))

		if exc.Name != c.Name && (!up.NameHidden && !up.NameReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Name", exc.Name, c.Name)
		}

		if exc.CampaignGroupId != c.CampaignGroupId {
			c.Comparisons = appendToChange(c.Comparisons, "Campaign Group ID", fmt.Sprintf("%v", exc.CampaignGroupId), fmt.Sprintf("%v", c.CampaignGroupId))
		}

		if exc.CampaignGroupWeight != c.CampaignGroupWeight {
			c.Comparisons = appendToChange(c.Comparisons, "Campaign Group Weight", fmt.Sprintf("%v", exc.CampaignGroupWeight), fmt.Sprintf("%v", c.CampaignGroupWeight))
		}

		if exc.AppendRc != c.AppendRc && (!up.AppendRcHidden && !up.AppendRcReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Append RC", fmt.Sprintf("%v", exc.AppendRc), fmt.Sprintf("%v", c.AppendRc))
		}

		if exc.Cpc != c.Cpc && (!up.CPCHidden && !up.CPCReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "CPC", fmt.Sprintf("$%v", exc.Cpc), fmt.Sprintf("$%v", c.Cpc))
		}

		if exc.Paused != c.Paused && (!up.PausedHidden && !up.PausedReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Paused", fmt.Sprintf("%v", exc.Paused), fmt.Sprintf("%v", c.Paused))
		}

		if exc.DailyImpsLimit != c.DailyImpsLimit && (!up.DailyBudgetHidden && !up.DailyBudgetReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Daily Imps Limit", fmt.Sprintf("%v", exc.DailyImpsLimit), fmt.Sprintf("%v", c.DailyImpsLimit))
		}

		if exc.TrackingMethod != c.TrackingMethod && (!up.TrackingMethodHidden && !up.TrackingMethodReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Tracking Method", fmt.Sprintf("%v", exc.TrackingMethod), fmt.Sprintf("%v", c.TrackingMethod))
		}

		if exc.Source != c.Source && (!up.SourceHidden && !up.SourceReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Source", fmt.Sprintf("%s", exc.Source), fmt.Sprintf("%s", c.Source))
		}

		if exc.PerformanceBasedPauseEnable != c.PerformanceBasedPauseEnable && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Performance Based Pause Enabled", fmt.Sprintf("%v", exc.PerformanceBasedPauseEnable), fmt.Sprintf("%v", c.PerformanceBasedPauseEnable))
		}

		if exc.PerformanceBasedCompareA != c.PerformanceBasedCompareA && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Performance Based Compare A", fmt.Sprintf("%v", exc.PerformanceBasedCompareA), fmt.Sprintf("%v", c.PerformanceBasedCompareA))
		}

		if exc.PerformanceBasedCompareB != c.PerformanceBasedCompareB && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Performance Based Compare B", fmt.Sprintf("%v", exc.PerformanceBasedCompareB), fmt.Sprintf("%v", c.PerformanceBasedCompareB))
		}

		if exc.PerformanceBasedPercent != c.PerformanceBasedPercent && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Performance Based Percent", fmt.Sprintf("%v", exc.PerformanceBasedPercent), fmt.Sprintf("%v", c.PerformanceBasedPercent))
		}

		if exc.PerformanceBasedNotifyOnly != c.PerformanceBasedNotifyOnly && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Performance Based Notify Only", fmt.Sprintf("%v", exc.PerformanceBasedNotifyOnly), fmt.Sprintf("%v", c.PerformanceBasedNotifyOnly))
		}

		if exc.EnableCampaignQualityCheck != c.EnableCampaignQualityCheck && (!up.EnablePerformanceBasedPauseHidden && !up.EnablePerformanceBasedPauseReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Enable Campaign Quality Check", fmt.Sprintf("%v", exc.EnableCampaignQualityCheck), fmt.Sprintf("%v", c.EnableCampaignQualityCheck))
		}

		if exc.StartDatetime.String() != c.StartDatetime.String() && (!up.EnableStartStopTimesHidden && !up.EnableStartStopTimesReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Start Date Time", fmt.Sprintf("%s", frt(exc.StartDatetime)), fmt.Sprintf("%s", frt(c.StartDatetime)))
		}

		if exc.EndDatetime.String() != c.EndDatetime.String() && (!up.EnableStartStopTimesHidden && !up.EnableStartStopTimesReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "End Date Time", fmt.Sprintf("%s", frt(exc.EndDatetime)), fmt.Sprintf("%s", frt(c.EndDatetime)))
		}

		if exc.DisableStartTime != c.DisableStartTime && (!up.EnableStartStopTimesHidden && !up.EnableStartStopTimesReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Disable Start Time", fmt.Sprintf("%v", exc.DisableStartTime), fmt.Sprintf("%v", c.DisableStartTime))
		}

		if exc.DisableEndTime != c.DisableEndTime && (!up.EnableStartStopTimesHidden && !up.EnableStartStopTimesReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Disable End Time", fmt.Sprintf("%v", exc.DisableEndTime), fmt.Sprintf("%v", c.DisableEndTime))
		}

		if exc.WeightVariance != c.WeightVariance && (!up.WeightVarianceHidden && !up.WeightVarianceReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Weight Variance", fmt.Sprintf("%v", exc.WeightVariance), fmt.Sprintf("%v", c.WeightVariance))
		}

		if exc.EnableStartStopTimes != c.EnableStartStopTimes && (!up.EnableStartStopTimesHidden && !up.EnableStartStopTimesReadOnly) {
			c.Comparisons = appendToChange(c.Comparisons, "Enable Start Stop Times", fmt.Sprintf("%v", exc.EnableStartStopTimes), fmt.Sprintf("%v", c.EnableStartStopTimes))
		}

		if !up.HideUrlUpdates {
			var curls []structs.AdscoopCampaignScheduleUrl
			var excurls []structs.AdscoopCampaignScheduleUrl
			db.Table("adscoop_scheduled_urls").Where("campaign_id = ?", c.ID).Find(&curls)
			db.Table("adscoop_urls").Where("campaign_id = ?", exc.ID).Find(&excurls)

			var varcurlss []string
			var excurlss []string

			for _, x := range curls {
				varcurlss = append(varcurlss, x.Url)
			}

			for _, x := range excurls {
				excurlss = append(excurlss, x.Url)
			}

			if fmt.Sprintf("%#s", varcurlss) != fmt.Sprintf("%#s", excurlss) {

				c.Comparisons = append(c.Comparisons, template.HTML("<tr><td colspan='3'><h6>New URL's</h6></td></tr>"))
				for _, x := range varcurlss {
					if !stringInSlice(x, excurlss) {
						c.Comparisons = append(c.Comparisons, template.HTML(fmt.Sprintf("<tr><td colspan='3'>%s</td></tr>", x)))
					}
				}

				c.Comparisons = append(c.Comparisons, template.HTML("<tr><td colspan='3'><h6>Removed URL's</h6></td></tr>"))
				for _, x := range excurlss {
					if !stringInSlice(x, varcurlss) {
						c.Comparisons = append(c.Comparisons, template.HTML(fmt.Sprintf("<tr><td colspan='3'>%s</td></tr>", x)))
					}
				}
			}
		}

		log.Printf("comparisons: %#s", c.Comparisons)

		retData.ScheduledCampaigns[i] = c

	}

	gun := mailgun.NewMailgun(config.MailGunDomain, config.MailGunApiKey, config.MailGunPublicApikey)

	if !user.IsBusinessHours() {
		retData.Client.ClientSchedulesPendApproval = true
	}

	if !retData.Client.ClientSchedulesPendApproval {
		gun := mailgun.NewMailgun(config.MailGunDomain, config.MailGunApiKey, config.MailGunPublicApikey)

		m := mailgun.NewMessage("donotnreply <donotreply@mg.adscoops.com>", fmt.Sprintf("Adscoops: Approved Client Updates for %s", retData.Client.Name), fmt.Sprintf("Client updates have been approved %s", retData.Client.Name))

		var output bytes.Buffer

		t, err := template.New("publishchangesforclient.tmpl").ParseFiles("templates/adscoops/publishchangesforclient.tmpl")

		if err != nil {
			log.Println("cannot load template: %s", err)
			return
		}

		err = t.Execute(&output, retData)

		if err != nil {
			log.Println("Template could not be executed: %s", err)
			return
		}

		var users []structs.AdscoopClientUser

		db.Where("client_id = ?", retData.Client.ID).Find(&users)

		m.SetHtml(output.String())

		for _, u := range users {
			m.AddRecipient(u.Email)
		}

		m.AddRecipient("daniel.aharonoff@broadscaler.com")
		m.AddRecipient("adops@adscoops.com")

		go func() {
			_, _, err = gun.Send(m)

			if err != nil {
				log.Println("Mailgun err: %s", err)
			}
		}()
	} else {
		m := mailgun.NewMessage("donotnreply <donotreply@mg.adscoops.com>", fmt.Sprintf("Adscoops: Client Updates for %s", retData.Client.Name), fmt.Sprintf("Client has made updates %s", retData.Client.Name))

		var output bytes.Buffer

		t, err := template.New("publishchanges.tmpl").ParseFiles("templates/adscoops/publishchanges.tmpl")

		if err != nil {
			log.Println("cannot load template: %s", err)
			return
		}

		err = t.Execute(&output, retData)

		if err != nil {
			log.Println("Template could not be executed: %s", err)
			return
		}

		m.SetHtml(output.String())

		m.AddRecipient("daniel.aharonoff@broadscaler.com")
		m.AddRecipient("adops@adscoops.com")

		go func() {
			_, _, err = gun.Send(m)
			if err != nil {
				log.Println("Mailgun err: %s", err)
			}
		}()
	}

}
