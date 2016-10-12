package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func broadvidVideosDomainsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.VideoDomains{}, c)
}

func broadvidVideosDomainsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.VideoDomain{}, c)
}

func broadvidVideosDomainsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.VideoDomain{}

	err = helpers.SaveEntity(ad, c, "Embed has been saved", ad.Host+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
