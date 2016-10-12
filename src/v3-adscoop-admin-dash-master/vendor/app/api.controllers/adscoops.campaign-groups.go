package apiControllers

import (

	"app/structs"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"app/helpers"
)

func campaignGroupssViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.CampaignGroups{}, c)
}

func campaignGroupsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.CampaignGroup{}, c)
}

func campaignGroupsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaignGroup = &structs.CampaignGroup{}
	err = helpers.SaveEntity(campaignGroup, c, "Campaign Group has been saved", "Campaign Group has been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
	}
}

func campaignGroupsByClient(c *gin.Context) {
	var clientCampaignGroups structs.ClientCampaignGroupReads

	err := clientCampaignGroups.Find(c.Param("id"))

	if err != nil {
		log.Errorf("Error getting client campaign groups: %s", err)
		c.JSON(500, err)
		return
	}

	c.JSON(200, clientCampaignGroups)
}

func campaignGroupsSaveByClientCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaignGroup = &structs.ClientCampaignGroupReads{}
	err = helpers.SaveEntityByRedir(campaignGroup, c, "Client Campaign Groups have been saved", "Client Campaign Groups have been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}
}

func campaignGroupsByRedirect(c *gin.Context) {
	var redirectCampaignGroups structs.RedirectCampaignGroupReads

	err := redirectCampaignGroups.Find(c.Param("id"))

	if err != nil {
		log.Errorf("Error getting redirect campaign groups: %s", err)
		c.JSON(500, err)
		return
	}

	c.JSON(200, redirectCampaignGroups)
}

func campaignGroupsSaveByRedirectCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var campaignGroup = &structs.RedirectCampaignGroupReads{}
	err = helpers.SaveEntityByRedir(campaignGroup, c, "Redirect campaign groups have been saved", "Redirect campaign groups have been saved", "success", 1, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}
}
