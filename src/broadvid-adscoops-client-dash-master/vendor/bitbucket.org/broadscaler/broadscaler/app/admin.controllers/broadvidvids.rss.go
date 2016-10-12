package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func broadvidVideosRsssViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.VideoRsss{}, c)
}

func broadvidVideosRsssViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.VideoRss{}, c)
}

func broadvidVideosRsssSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.VideoRss{}

	err = helpers.SaveEntity(ad, c, "Embed has been saved", ad.Label+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
