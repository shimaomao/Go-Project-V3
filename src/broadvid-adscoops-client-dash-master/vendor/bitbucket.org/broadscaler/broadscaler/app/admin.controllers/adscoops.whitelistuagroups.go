package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func adscoopsWhitelistUaGroupsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.AdscoopWhitelistUseragentGroups{}, c)
}

func adscoopsWhitelistUaGroupsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.AdscoopWhitelistUseragentGroup{}, c)
}

func adscoopsWhitelistUaGroupsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.AdscoopWhitelistUseragentGroup{}

	err = helpers.SaveEntity(ad, c, "Whitelist URL Group has been saved", ad.Name+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
