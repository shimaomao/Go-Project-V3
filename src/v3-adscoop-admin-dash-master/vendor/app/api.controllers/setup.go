package apiControllers

import (
	"app/adonnetwork"
	"app/adscoops.caches"
	"app/helpers"
	log "github.com/Sirupsen/logrus"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

/*Setup is to set up the controllers for the admin for adscoops and broadvids*/
func Setup(m *gin.Engine) {
	log.Println("Setting up admin controllers...")

	helpers.AdscoopsDB = AdscoopsDB
	helpers.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	helpers.BroadvidDB = BroadvidDB
	helpers.RedisPool = RedisPool
	helpers.LinkToStructs()

	adscoopsCaches.AdscoopsDB = AdscoopsDB
	adscoopsCaches.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	adscoopsCaches.BroadvidDB = BroadvidDB
	adscoopsCaches.RedisPool = RedisPool

	adonnetwork.AdscoopsDB = AdscoopsDB

	m.GET("/auth/callback", AuthCallbackCtrl)

	auth := m.Group("/", RequireAuth)

	m.GET("/auth/login", func(c *gin.Context) {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	m.GET("/auth/logout", func(c *gin.Context) {
		ses := ginsessions.Default(c)

		ses.Delete("UserID")
		ses.Save()
		c.Redirect(302, c.Request.Referer())
	})

	m.GET("/user/info", UserInfoCtrl)

	asc := m.Group("/adscoops")
	{
		ascadon := asc.Group("/advertiser")
		{
			ascadon.GET("/campaigns/viewall", adonNetworkCampaignsViewAllCtrl)
		}

		ascclients := asc.Group("/clients")
		{
			ascclients.GET("/viewall", clientsViewallCtrl)
			ascclients.GET("/viewVisible", clientsViewVisibleCtrl)
			ascclients.GET("/view/:id", clientsViewCtrl)
			ascclients.POST("/save", clientsSaveCtrl)
			ascclients.POST("/showReport/:id", clientsShowReportCtrl)
			ascclients.POST("/updateCampaignSort", clientsUpdateCampaignSortPerUserCtrl)

			// Payment endpoints
			ascclients.GET("/newPaymentHash/:id", clientsNewPaymentHashCtrl)
			ascclients.POST("/manualCharge/:id", clientsManualChargeCtrl)
			ascclients.POST("/autoCharge/:id", clientsAutoChargeCtrl)

			// stats

			ascclients.GET("/viewRedirStats/:id", clientsViewRedirStatsCtrl)

			// redirects
			ascclients.GET("/viewAssociatedRedirects/:id", clientsViewAssociatedRedirectsCtrl)
		}
		asccampaigns := asc.Group("/campaigns")
		{
			asccampaigns.GET("/viewall", campaignsViewallCtrl)
			asccampaigns.GET("/view/:id", campaignsViewCtrl)
			asccampaigns.GET("/byRedirect/:id", campaignsByRedirectCtrl)
			asccampaigns.POST("/saveByRedirect/:id", campaignsSaveByRedirectCtrl)
			asccampaigns.POST("/basicSave", campaignsBasicSaveCtrl)
			asccampaigns.POST("/save", campaignsSaveCtrl)
			asccampaigns.GET("/clientviewall/:id", campaignsClientViewAllCtrl)
			asccampaigns.GET("/client-viewextradetails/:id", campaignsClientViewExtraDetailsCtrl)
		}

		asccampaigngroups := asc.Group("/campaign-groups")
		{
			asccampaigngroups.GET("/viewall", campaignGroupssViewallCtrl)
			asccampaigngroups.GET("/view/:id", campaignGroupsViewCtrl)
			asccampaigngroups.POST("/save", campaignGroupsSaveCtrl)

			asccampaigngroups.GET("/byClient/:id", campaignGroupsByClient)
			asccampaigngroups.POST("/saveByClient/:id", campaignGroupsSaveByClientCtrl)

			asccampaigngroups.GET("/byRedirect/:id", campaignGroupsByRedirect)
			asccampaigngroups.POST("/saveByRedirect/:id", campaignGroupsSaveByRedirectCtrl)

		}

		asccampaignschedules := asc.Group("/campaign-schedules")
		{
			asccampaignschedules.GET("/viewall/:id", campaignSchedulesViewallCtrl)
			asccampaignschedules.GET("/view/:id", campaignSchedulesViewCtrl)
			asccampaignschedules.POST("/save", campaignSchedulesSaveCtrl)
		}

		ascredirects := asc.Group("/redirects")
		{
			ascredirects.GET("/viewall", redirectsViewallCtrl)
			ascredirects.GET("/view/:id", redirectsViewCtrl)
			ascredirects.POST("/save", redirectsSaveCtrl)
		}

		aschosts := asc.Group("/hosts")
		{
			aschosts.GET("/viewall", hostsViewallCtrl)
		}

		ascwhitelisturlgroups := asc.Group("/whitelisturlgroups")
		{
			ascwhitelisturlgroups.GET("/viewall", adscoopsWhitelistUrlGroupsViewallCtrl)
			ascwhitelisturlgroups.GET("/view/:id", adscoopsWhitelistUrlGroupsViewCtrl)
			ascwhitelisturlgroups.POST("/save", adscoopsWhitelistUrlGroupsSaveCtrl)
		}

		ascwhitelistuagroups := asc.Group("/whitelistuagroups")
		{
			ascwhitelistuagroups.GET("/viewall", adscoopsWhitelistUaGroupsViewallCtrl)
			ascwhitelistuagroups.GET("/view/:id", adscoopsWhitelistUaGroupsViewCtrl)
			ascwhitelistuagroups.POST("/save", adscoopsWhitelistUaGroupsSaveCtrl)
		}

		ascstats := asc.Group("/stats")
		{
			ascstats.GET("/dailyImpressionCount", statsDailyImpressionCountCtrl)

			ascstats.GET("/getVerticalStats", statsDailyImpressionCountByVerticalCtrl)

			ascstats.GET("/getRealTimeStats/today", statsRealTimeTodayCtrl)
			ascstats.GET("/getRealTimeStats/yesterday", statsRealTimeYesterdayCtrl)
			ascstats.GET("/getRealTimeClientStats", statsRealTimeClientCtrl)
		}
	}

	bvads := m.Group("/broadvidads")
	{

		bvadsads := bvads.Group("/ads")
		{
			bvadsads.GET("/viewall", broadvidAdsViewallCtrl)
			bvadsads.GET("/view/:id", broadvidAdsViewCtrl)
			bvadsads.POST("/save", broadvidAdsSaveCtrl)
		}

		bvadsadsembeds := bvads.Group("/ad_embeds")
		{
			bvadsadsembeds.GET("/viewall", broadvidAdEmbedsViewallCtrl)
			bvadsadsembeds.GET("/view/:id", broadvidAdEmbedsViewCtrl)
			bvadsadsembeds.POST("/save", broadvidAdEmbedsSaveCtrl)
			bvadsadsembeds.POST("/ad_embeds/remove", broadvidadsAdEmbedsRemoveCtrl)
			bvadsadsembeds.POST("/ad_embeds/copy", broadvidadsAdEmbedsCopyCtrl)
			bvadsadsembeds.POST("/ad_embeds/pause", broadvidadsAdEmbedsPauseCtrl)
		}

		bvadsembeds := bvads.Group("/embeds")
		{
			bvadsembeds.GET("/viewall", broadvidEmbedsViewallCtrl)
			bvadsembeds.GET("/view/:id", broadvidEmbedsViewCtrl)
			bvadsembeds.POST("/save", broadvidEmbedsSaveCtrl)
		}

		bvwlug := bvads.Group("/whitelisturlgroups")
		{
			bvwlug.GET("/viewall", broadvidWhitelistUrlGroupsViewallCtrl)
			bvwlug.GET("/view/:id", broadvidWhitelistUrlGroupsViewCtrl)
			bvwlug.POST("/save", broadvidWhitelistUrlGroupsSaveCtrl)
		}

		bvblug := bvads.Group("/blacklisturlgroups")
		{
			bvblug.GET("/viewall", broadvidBlacklistUrlGroupsViewallCtrl)
			bvblug.GET("/view/:id", broadvidBlacklistUrlGroupsViewCtrl)
			bvblug.POST("/save", broadvidBlacklistUrlGroupsSaveCtrl)
		}

		bvwluag := bvads.Group("/whitelistuagroups")
		{
			bvwluag.GET("/viewall", broadvidWhitelistUaGroupsViewallCtrl)
			bvwluag.GET("/view/:id", broadvidWhitelistUaGroupsViewCtrl)
			bvwluag.POST("/save", broadvidWhitelistUaGroupsSaveCtrl)
		}

		bvwlcg := bvads.Group("/whitelistcountrygroups")
		{
			bvwlcg.GET("/viewall", broadvidWhitelistCountryGroupsViewallCtrl)
			bvwlcg.GET("/view/:id", broadvidWhitelistCountryGroupsViewCtrl)
			bvwlcg.POST("/save", broadvidWhitelistCountryGroupsSaveCtrl)
		}

		bvads.POST("/flushadcache", broadvidadsFlushAdCache)
	}

	bvvids := m.Group("/broadvidvideos")
	{

		bvvidsrss := bvvids.Group("/rss")
		{
			bvvidsrss.GET("/viewall", broadvidVideosRsssViewallCtrl)
			bvvidsrss.GET("/view/:id", broadvidVideosRsssViewCtrl)
			bvvidsrss.POST("/save", broadvidVideosRsssSaveCtrl)
		}

		bvvidsembeds := bvvids.Group("/embeds")
		{
			bvvidsembeds.GET("/viewall", broadvidVideosEmbedsViewallCtrl)
			bvvidsembeds.GET("/view/:id", broadvidVideosEmbedsViewCtrl)
			bvvidsembeds.POST("/save", broadvidVideosEmbedsSaveCtrl)
		}

		bvvidsredirects := bvvids.Group("/redirects")
		{
			bvvidsredirects.GET("/viewall", broadvidVideosRedirectsViewallCtrl)
			bvvidsredirects.GET("/view/:id", broadvidVideosRedirectsViewCtrl)
			bvvidsredirects.POST("/save", broadvidVideosRedirectsSaveCtrl)
		}

		bvvidsinjectjs := bvvids.Group("/injectjs")
		{
			bvvidsinjectjs.GET("/viewall", broadvidVideosInjectJssViewallCtrl)
			bvvidsinjectjs.GET("/view/:id", broadvidVideosInjectJssViewCtrl)
			bvvidsinjectjs.POST("/save", broadvidVideosInjectJssSaveCtrl)
		}

		bvvidsthemes := bvvids.Group("/themes")
		{
			bvvidsthemes.GET("/viewall", broadvidVideosThemesViewallCtrl)
			bvvidsthemes.GET("/view/:id", broadvidVideosThemesViewCtrl)
			bvvidsthemes.POST("/save", broadvidVideosThemesSaveCtrl)
		}

		bvvidsdomains := bvvids.Group("/domains")
		{
			bvvidsdomains.GET("/viewall", broadvidVideosDomainsViewallCtrl)
			bvvidsdomains.GET("/view/:id", broadvidVideosDomainsViewCtrl)
			bvvidsdomains.POST("/save", broadvidVideosDomainsSaveCtrl)
		}
	}

	settings := m.Group("/settings")
	{
		settingsusers := settings.Group("/users")
		{
			settingsusers.GET("/viewall", settingsUsersViewallCtrl)
			settingsusers.GET("/view", settingsUsersViewCtrl)
			settingsusers.POST("/save", settingsUsersSaveCtrl)
		}
	}

	auth.StaticFile("/", "./public/index.html")
	auth.StaticFile("/bundle.js", "./public/bundle.js")
	auth.Static("/assets", "./public/assets")
	auth.Static("/node_modules", "./public/node_modules")
	auth.Static("/custom", "./public/custom")
	auth.Static("/bower_components", "./public/bower_components")
	auth.Static("/partials", "./public/partials")
	auth.Static("/lib", "./public/lib")
	auth.Static("/app", "./public/app")
	auth.Static("/views", "./public/views")

	auth.GET("/ping", PingCtrl)

	auth.GET("/rtupdates", RealtimeUpdatesCtrl)

	auth.GET("/adscoopsupdates", AdscoopsRealtimeUpdatesCtrl)

	go statsUpdates()

}
