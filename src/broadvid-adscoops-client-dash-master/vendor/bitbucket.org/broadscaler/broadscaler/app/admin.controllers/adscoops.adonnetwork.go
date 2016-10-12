package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/adonnetwork"
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"github.com/gin-gonic/gin"
)

func adonNetworkCampaignsViewAllCtrl(c *gin.Context) {
	helpers.FindAll(&adonnetwork.Campaigns{}, c)
}
