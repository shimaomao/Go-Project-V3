package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func broadvidWhitelistCountryGroupsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.AdCountryGroups{}, c)
}

func broadvidWhitelistCountryGroupsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.AdCountryGroup{}, c)
}

func broadvidWhitelistCountryGroupsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.AdCountryGroup{}

	err = helpers.SaveEntity(ad, c, "Embed has been saved", ad.Name+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
