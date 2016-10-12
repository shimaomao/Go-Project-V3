package apiControllers

import (
	"adonnetwork"
	"helpers"
	"github.com/gin-gonic/gin"
)

func adonNetworkCampaignsViewAllCtrl(c *gin.Context) {
	helpers.FindAll(&adonnetwork.Campaigns{}, c)
}
