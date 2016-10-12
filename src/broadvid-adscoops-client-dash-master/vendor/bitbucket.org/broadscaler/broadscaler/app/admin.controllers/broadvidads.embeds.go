package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/helpers"
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func broadvidEmbedsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.Embeds{}, c)
}

func broadvidEmbedsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.Embed{}, c)
}

func broadvidEmbedsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.Embed{}

	err = helpers.SaveEntity(ad, c, "Embed has been saved", ad.Label+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
