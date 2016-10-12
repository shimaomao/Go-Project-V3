package apiControllers

import (
	log "github.com/Sirupsen/logrus"

	"app/adscoops.caches"
	"github.com/gin-gonic/gin"
	"app/helpers"
	"app/structs"
)

func campaignsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.Campaigns{}, c)
}

func campaignsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.Campaign{}, c)
}

func campaignsByRedirectCtrl(c *gin.Context) {
	var redirectCampaigns structs.RedirectCampaignReads

	err := redirectCampaigns.Find(c.Param("id"))

	if err != nil {
		log.Errorf("Error getting redirect campaigns: %s", err)
		c.JSON(500, err)
		return
	}

	c.JSON(200, redirectCampaigns)
}

func campaignsSaveByRedirectCtrl(c *gin.Context) {

	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaign = &structs.RedirectCampaignReads{}
	err = helpers.SaveEntityByRedir(campaign, c, "Redirect Campaigns have been saved", "Redirect Campaigns have been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}
}

func campaignsClientViewAllCtrl(c *gin.Context) {
	var campaigns structs.Campaigns
	campaigns.FindFromClient(c.Param("id"))
	c.JSON(200, campaigns)
}
func campaignsClientViewExtraDetailsCtrl(c *gin.Context) {
	var campaigns structs.Campaigns
	campaigns.FindExtraDetailsFromClient(c.Param("id"))
	c.JSON(200, campaigns)
}

func campaignsSaveCtrl(c *gin.Context) {

	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaign = &structs.Campaign{}
	err = helpers.SaveEntity(campaign, c, "Campaign has been saved", campaign.Name+" has been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

	adscoopsCaches.UpdateCampaign(campaign.ID)
}

func campaignsBasicSaveCtrl(c *gin.Context) {

	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaign = &structs.Campaign{}
	err = helpers.BasicSaveEntity(campaign, c, "Campaign has been saved", campaign.Name+" has been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

	adscoopsCaches.UpdateCampaign(campaign.ID)
}
