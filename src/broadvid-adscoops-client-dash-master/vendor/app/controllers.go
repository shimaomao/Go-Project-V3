package app

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func frt(dateTime time.Time) (ret string) {

	ret = dateTime.Format("01/02/2006 03:04 PM")
	return
}

func frtLA(dateTime time.Time) (ret string) {
	pst := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
		dateTime.Hour(), dateTime.Minute(), 0, 0, time.UTC)
	location, _ := time.LoadLocation("America/Los_Angeles")

	dateTime = pst.In(location)
	ret = dateTime.Format("01/02/2006 03:04 PM")
	return
}

func appendToChange(c []template.HTML, field string, was string, is string) []template.HTML {
	c = append(c, template.HTML(fmt.Sprintf("<tr style='border-bottom:1px solid #000'><td style='background-color:#ddd; padding:1em;'>%s</td><td style='background-color:#eee; padding:1em;'>%s</td><td style='padding:1em;'>%s</td></tr>", field, was, is)))
	return c
}

func controllersSetup(m *martini.ClassicMartini) {

	m.Get("/", requireLogin, func(r render.Render, user *UserWithPolicy, req *http.Request) {
		var retData struct {
			DefaultRetData
			Timestamp          time.Time
			EditedSchedules    bool
			Campaigns          CacheCampaign
			UserPolicy         structs.AdscoopUserPolicy
			ScheduledCampaigns []struct {
				structs.AdscoopCampaign
				PendingSchedule uint      `gorm:"column:asc_id"`
				ScheduleLabel   string    `gorm:"column:schedule_label"`
				ScheduleExecute time.Time `gorm:"column:schedule_execute"`
			}
		}
		retData.User = user

		retData.Campaigns = campaignsCache[user.ClientID]

		retData.Timestamp = time.Now()

		db.Select("adscoop_campaigns.*, adscoop_campaign_schedules.id AS asc_id, schedule_label, schedule_execute").Table("adscoop_campaigns").Joins("LEFT JOIN adscoop_campaign_schedules ON adscoop_campaign_schedules.campaign_id = adscoop_campaigns.id AND adscoop_campaign_schedules.schedule_pending != 0").Where("adscoop_campaign_schedules.schedule_pending > 0 AND adscoop_campaigns.client_id = ?", user.ClientID).Find(&retData.ScheduledCampaigns)

		for _, x := range retData.ScheduledCampaigns {
			if x.PendingSchedule != 0 {
				retData.EditedSchedules = true
			}
		}

		db.Find(&retData.UserPolicy, user.UserPolicyID)
		if req.FormValue("ajax") == "1" {
			r.HTML(http.StatusOK, "adscoops/clients/campaigns/campaignList", retData, render.HTMLOptions{})
		} else {
			r.HTML(http.StatusOK, "adscoops/clients/campaigns/list", retData)
		}
	})

	m.Get("/url-encoding", requireLogin, func(r render.Render, user *UserWithPolicy) {

		var retData struct {
			DefaultRetData
			Urls          []string
			FormRedirPage string
			FormUrls      string
		}
		retData.User = user
		r.HTML(http.StatusOK, "adscoops/urlencoding", retData)
	})

	m.Post("/url-encoding", requireLogin, func(req *http.Request, r render.Render, user *UserWithPolicy) {
		redirPage, urls := req.FormValue("redir_page"), req.FormValue("urls")

		urlsList := strings.SplitN(urls, "\r\n", -1)

		retUrls := []string{}

		for _, u := range urlsList {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			if strings.HasPrefix(u, "Campaign: ") {
				retUrls = append(retUrls, "")
				retUrls = append(retUrls, u)
				continue
			}
			if !strings.HasPrefix(u, "http") {
				u = "http://" + u
			}
			extu, _ := url.Parse(u)

			newDomain := url.URL{}
			newDomain.Scheme = extu.Scheme
			newDomain.Host = extu.Host
			q := newDomain.Query()
			q.Set("url", extu.String())
			newDomain.RawQuery = q.Encode()
			newDomain.Path = redirPage

			retUrls = append(retUrls, newDomain.String())
		}

		var retData struct {
			DefaultRetData
			Urls          []string
			FormRedirPage string
			FormUrls      string
		}

		retData.Urls = retUrls
		retData.FormUrls = urls
		retData.FormRedirPage = redirPage
		retData.User = user

		r.HTML(http.StatusOK, "adscoops/urlencoding", retData)
	})

	m.Get("/login", func(r render.Render) {
		var retData struct {
			DefaultRetData
		}

		retData.User = &UserWithPolicy{}
		r.HTML(http.StatusOK, "login", retData)
	})

	m.Post("/login", binding.Bind(structs.AdscoopClientUser{}), func(r render.Render, user structs.AdscoopClientUser, s sessions.Session) {
		loginError := "User and/or password does not match or user is not found"
		var tempUser structs.AdscoopClientUser
		db.Where("email = ?", user.Email).Find(&tempUser)

		if tempUser.ID == 0 {
			r.Text(http.StatusInternalServerError, loginError)
			return
		}

		password := []byte(user.Password)
		hashedPassword := []byte(tempUser.Password)

		err := bcrypt.CompareHashAndPassword(hashedPassword, password)

		if err != nil {
			r.Text(http.StatusInternalServerError, loginError)
			return
		}

		var client structs.AdscoopClient
		db.Where("id = ? AND enable_client_login = 1", tempUser.ClientID).Find(&client)

		if tempUser.ID == 0 {
			r.Text(http.StatusInternalServerError, loginError)
			return
		}

		s.Set("user_id", tempUser.ID)

		go createCampaignCache(tempUser.ClientID)
		r.Redirect("/")
	})

	m.Get("/publish-updates", requireLogin, func(r render.Render, user *UserWithPolicy) {
		var schedules []structs.AdscoopCampaignSchedule
		db.Select("adscoop_campaign_schedules.*").Table("adscoop_campaign_schedules").Joins("JOIN adscoop_campaigns ON adscoop_campaign_schedules.campaign_id = adscoop_campaigns.id").Where("adscoop_campaign_schedules.schedule_pending != 0 AND adscoop_campaigns.client_id = ?", user.ClientID).Find(&schedules)

		var client structs.AdscoopClient

		db.Find(&client, user.ClientID)

		sendEmailForClient(user)

		now := time.Now().UTC()

		for _, sch := range schedules {
			if user.IsBusinessHours() {
				sch.ScheduleQueued = client.ClientSchedulesPendApproval
			} else {
				sch.ScheduleQueued = true
			}
			sch.SchedulePending = false
			if user.Policy.AutoApproveDelay != 0 && sch.ScheduleQueued {
				delayTime := now.Add(time.Duration(user.Policy.AutoApproveDelay) * time.Minute)
				if sch.ScheduleExecute.Before(delayTime) {
					sch.ScheduleExecute = delayTime
				}
			}
			db.Save(&sch)
		}

		r.Redirect("/")
	})

	m.Post("/logout", func(r render.Render, s sessions.Session) {
		ops := sessions.Options{
			MaxAge: -1,
		}
		s.Options(ops)
		s.Clear()
		r.Redirect("/login")
	})

	m.Get("/logout", requireLogin, func(r render.Render, s sessions.Session, user *UserWithPolicy) {
		var sc []structs.AdscoopCampaign
		db.Select("*").Table("adscoop_campaigns").Joins("LEFT JOIN adscoop_campaign_schedules ON adscoop_campaign_schedules.campaign_id = adscoop_campaigns.id AND adscoop_campaign_schedules.schedule_pending != 0").Where("adscoop_campaign_schedules.schedule_pending > 0 AND adscoop_campaigns.client_id = ?", user.ClientID).Find(&sc)

		if len(sc) != 0 {
			r.HTML(200, "confirm_logout", nil, render.HTMLOptions{})
			return
		}

		ops := sessions.Options{
			MaxAge: -1,
		}
		s.Options(ops)
		s.Clear()
		r.Redirect("/login")
	})

	m.Get("/campaigns/new-campaign", requireLogin, adscoopsClientsCampaignsScheduleNew)
	m.Get("/campaigns/:cid/new-schedule", requireLogin, adscoopsClientsCampaignsScheduleNew)
	m.Get("/campaigns/:cid/edit-schedule/:sid", requireLogin, adscoopsClientsCampaignsScheduleNew)
	m.Post("/campaigns/:cid/save", binding.Bind(structs.AdscoopCampaignSchedule{}), requireLogin, adscoopsClientsCampaignsScheduleSave)
}
