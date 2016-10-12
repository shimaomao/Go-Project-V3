package apiControllers

import (
	"app/adonnetwork"
	"github.com/gin-gonic/gin"
	"app/helpers"
)

func adonNetworkCampaignsViewAllCtrl(c *gin.Context) {
	helpers.FindAll(&adonnetwork.Campaigns{}, c)
}
