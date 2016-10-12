package apiControllers

import (
	"app/adonnetwork"
	"app/helpers"
	"github.com/gin-gonic/gin"
)

func adonNetworkCampaignsViewAllCtrl(c *gin.Context) {
	helpers.FindAll(&adonnetwork.Campaigns{}, c)
}
