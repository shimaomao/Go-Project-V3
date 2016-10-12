package apiControllers

import (
	"github.com/gin-gonic/gin"
	"app/helpers"
	"app/structs"
)

func broadvidBlacklistUrlGroupsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.AdUrlGroupBlacklists{}, c)
}

func broadvidBlacklistUrlGroupsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.AdUrlGroupBlacklist{}, c)
}

func broadvidBlacklistUrlGroupsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.AdUrlGroupBlacklist{}

	err = helpers.SaveEntity(ad, c, "Blacklist Url Group has been saved", ad.Name+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
