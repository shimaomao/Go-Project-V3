package apiControllers

import (
	"app/helpers"
	"app/structs"
	"github.com/gin-gonic/gin"
)

func adscoopsWhitelistUrlGroupsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.AdscoopWhitelists{}, c)
}

func adscoopsWhitelistUrlGroupsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.AdscoopWhitelist{}, c)
}

func adscoopsWhitelistUrlGroupsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.AdscoopWhitelist{}

	err = helpers.SaveEntity(ad, c, "Whitelist URL Group has been saved", ad.Name+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
