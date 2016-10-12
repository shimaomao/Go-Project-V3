package apiControllers

import (
	"adonnetwork"
	"github.com/gin-gonic/gin"
	"helpers"
)

func adonNetworkCampaignsViewAllCtrl(c *gin.Context) {
	helpers.FindAll(&adonnetwork.Campaigns{}, c)
}
